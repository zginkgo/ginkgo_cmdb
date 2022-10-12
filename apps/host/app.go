package host

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
	"github.com/zginkgo/ginkgo_cmdb/apps/resource"
	"github.com/zginkgo/ginkgo_cmdb/utils"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	AppName = "host"
	// DefaultPageSize 默认分页大小
	DefaultPageSize = 20
	// DefaultPageNumber 默认页号
	DefaultPageNumber = 1
)

var (
	validate = validator.New()
)

func NewDefaultHost() *Host {
	return &Host{
		Base: &resource.Base{
			ResourceType: resource.Type_HOST,
		},
		Information: &resource.Information{},
		Describe:    &Describe{},
	}
}

func (h *Host) GenHash() error {
	// hash resource
	h.Base.ResourceHash = h.Information.Hash()

	// hash describe
	h.Base.DescribeHash = utils.Hash(h.Describe)
	return nil
}

func (d *Describe) KeyPairNameToString() string {
	return strings.Join(d.KeyPairName, ",")
}

func (d *Describe) SecurityGroupsToString() string {
	return strings.Join(d.SecurityGroups, ",")
}

func (d *Describe) LoadKeyPairNameString(s string) {
	if s != "" {
		d.KeyPairName = strings.Split(s, ",")
	}
}

func (d *Describe) LoadSecurityGroupsString(s string) {
	if s != "" {
		d.SecurityGroups = strings.Split(s, ",")
	}
}

func (req *DescribeHostRequest) Where() (string, interface{}) {
	switch req.DescribeBy {
	default:
		return "r.id = ?", req.Value
	}
}

func NewDescribeHostRequestWithID(id string) *DescribeHostRequest {
	return &DescribeHostRequest{
		DescribeBy: DescribeBy_HOST_ID,
		Value:      id,
	}
}

func NewDeleteHostRequestWithID(id string) *ReleaseHostRequest {
	return &ReleaseHostRequest{Id: id}
}

func NewUpdateHostRequest(id string) *UpdateHostRequest {
	return &UpdateHostRequest{
		Id:             id,
		UpdateMode:     UpdateMode_PUT,
		UpdateHostData: &UpdateHostData{},
	}
}

func (req *UpdateHostRequest) Validate() error {
	return validate.Struct(req)
}

func (h *Host) Put(req *UpdateHostData) {
	oldRH, oldDH := h.Base.ResourceHash, h.Base.DescribeHash

	h.Information = req.Information
	h.Describe = req.Describe
	h.Information.UpdateAt = time.Now().UnixMicro()
	h.GenHash()

	if h.Base.ResourceHash != oldRH {
		h.Base.ResourceHashChanged = true
	}
	if h.Base.DescribeHash != oldDH {
		h.Base.DescribeHashChanged = true
	}
}

func (h *Host) ShortDesc() string {
	return fmt.Sprintf("%s %s", h.Information.Name, h.Information.PrivateIp)
}

func NewUpdateHostDataByIns(ins *Host) *UpdateHostData {
	return &UpdateHostData{
		Information: ins.Information,
		Describe:    ins.Describe,
	}
}

func NewHostSet() *HostSet {
	return &HostSet{
		Items: []*Host{},
	}
}

func (s *HostSet) Add(item any) {
	s.Items = append(s.Items, item.(*Host))
	return
}

func (s *HostSet) Length() int64 {
	return int64(len(s.Items))
}

func (s *HostSet) Clone() *HostSet {
	return proto.Clone(s).(*HostSet)
}

func (s *HostSet) ResourceIds() (ids []string) {
	for i := range s.Items {
		ids = append(ids, s.Items[i].Base.Id)
	}
	return
}

func (s *HostSet) UpdateTag(tags []*resource.Tag) {
	for i := range tags {
		for j := range s.Items {
			if s.Items[j].Base.Id == tags[i].ResourceId {
				s.Items[j].Information.AddTag(tags[i])
			}
		}
	}
}

// NewPageRequestFromHTTP 从HTTP请求中加载分页请求
func NewPageRequestFromHTTP(req *http.Request) *PageRequest {
	qs := req.URL.Query()

	ps := qs.Get("page_size")
	pn := qs.Get("page_number")
	os := qs.Get("offset")

	psUint64, _ := strconv.ParseUint(ps, 10, 64)
	pnUint64, _ := strconv.ParseUint(pn, 10, 64)
	osInt64, _ := strconv.ParseInt(os, 10, 64)

	if psUint64 == 0 {
		psUint64 = DefaultPageSize
	}
	if pnUint64 == 0 {
		pnUint64 = DefaultPageNumber
	}

	return &PageRequest{
		PageSize:   psUint64,
		PageNumber: pnUint64,
		Offset:     osInt64,
	}
}

func NewQueryHostRequestFromHTTP(r *http.Request) *QueryHostRequest {
	qs := r.URL.Query()
	page := NewPageRequestFromHTTP(r)
	kw := qs.Get("keywords")

	return &QueryHostRequest{
		Page:     page,
		Keywords: kw,
	}
}

// 分页器
// for p.Next() {
// 	if err := p.Scan(set); err != nil {
// 		...
// 	}
// }
type Pagger interface {
	Next() bool
	SetPageSize(ps int64)
	Scan(context.Context, *HostSet) error
}

func NewDefaultPageRequest() *PageRequest {
	return NewPageRequest(DefaultPageSize, DefaultPageNumber)
}

// NewPageRequest 实例化
func NewPageRequest(ps uint, pn uint) *PageRequest {
	return &PageRequest{
		PageSize:   uint64(ps),
		PageNumber: uint64(pn),
	}
}

// GetOffset skip
// 如果传入了offset则使用传入的offset参数
func (p *PageRequest) ComputeOffset() int64 {
	if p.Offset != 0 {
		return p.Offset
	}

	return int64(p.PageSize * (p.PageNumber - 1))
}

// 抽象调用Pager
type Set interface {
	// 往Set里面添加元素, 任何类型都可以
	Add(any)

	// 当前的集合有多个元素
	Length() int64
}

type PagerV2 interface {
	Next() bool
	Scan(context.Context, Set) error
	Offset() int64
	SetPageSize(ps int64)
	SetRate(r float64)
	PageSize() int64
	PageNumber() int64
}

func NewBasePagerV2() *BasePagerV2 {
	return &BasePagerV2{
		hasNext:    true,
		tb:         tokenbucket.NewBucketWithRate(1, 1),
		pageNumber: 1,
		pageSize:   20,
	}
}

// 面向组合, 用他来实现一个模板, 除了Scan的其他方法都实现

type BasePagerV2 struct {
	// 令牌桶
	hasNext bool
	tb      *tokenbucket.Bucket

	// 控制分页的核心参数
	pageNumber int64
	pageSize   int64
}

func (p *BasePagerV2) Next() bool {
	// 等待分配令牌
	p.tb.Wait(1)

	return p.hasNext
}

func (p *BasePagerV2) Offset() int64 {
	return (p.pageNumber - 1) * p.pageSize
}

func (p *BasePagerV2) SetPageSize(ps int64) {
	p.pageSize = ps
}

func (p *BasePagerV2) PageSize() int64 {
	return p.pageSize
}

func (p *BasePagerV2) PageNumber() int64 {
	return p.pageNumber
}

func (p *BasePagerV2) SetRate(r float64) {
	p.tb.SetRate(r)
}

func (p *BasePagerV2) CheckHasNext(current int64) {
	// 可以根据当前一页是满页来决定是否有下一页
	if current < p.pageSize {
		p.hasNext = false
	} else {
		// 直接调整指针到下一页
		p.pageNumber++
	}
}
