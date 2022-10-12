package resource

import (
	"fmt"
	"github.com/zginkgo/ginkgo_cmdb/utils"
	"net/http"
	"strconv"
	"strings"
)

const (
	AppName = "resource"
)

const (
	// DefaultPageSize 默认分页大小
	DefaultPageSize = 20
	// DefaultPageNumber 默认页号
	DefaultPageNumber = 1
)

// Operator Tag的比较操作符, 内比promethues 的Tag比较操作, 官网也能找到该4种操作符号
type Operator string

const (
	// Operator_EQUAL SQL 比较操作  =
	Operator_EQUAL = "="

	// Operator_NOT_EQUAL SQL 比较操作  !=
	Operator_NOT_EQUAL = "!="

	// SQL 比较操作  LIKE
	Operator_LIKE_EQUAL = "=~"
	// SQL 比较操作  NOT LIKE
	Operator_NOT_LIKE_EQUAL = "!~"
)

func (r *SearchRequest) HasTag() bool {
	return len(r.Tags) > 0
}

// 多个值比较的关系说明:
//    app=~app1,app2  你不能说 app1和app2是 AND关系, 一定是OR关系    是一种白名单策略(包含策略)
//    app!~app3,app4  tag_key=app tag_value NOT LIKE (app3,app4), 是一种黑名单策略(排除策略)
func (s *TagSelector) RelationShip() string {
	switch s.Operator {
	case Operator_EQUAL, Operator_LIKE_EQUAL:
		return " OR "
	case Operator_NOT_EQUAL, Operator_NOT_LIKE_EQUAL:
		return " AND "
	default:
		return " OR "
	}
}

func NewResourceSet() *ResourceSet {
	return &ResourceSet{
		Items: []*Resource{},
	}
}

func (s *ResourceSet) Add(item *Resource) {
	s.Items = append(s.Items, item)
}

func NewDefaultResource() *Resource {
	return &Resource{
		Base:        &Base{},
		Information: &Information{},
	}
}

func (i *Information) LoadPrivateIPString(s string) {
	if s != "" {
		i.PrivateIp = strings.Split(s, ",")
	}
}

func (i *Information) LoadPublicIPString(s string) {
	if s != "" {
		i.PublicIp = strings.Split(s, ",")
	}
}

func NewDefaultTag() *Tag {
	return &Tag{
		Type:   TagType_USER,
		Weight: 1,
	}
}

func (s *ResourceSet) ResourceIds() (ids []string) {
	for i := range s.Items {
		ids = append(ids, s.Items[i].Base.Id)
	}

	return
}

func (s *ResourceSet) UpdateTag(tags []*Tag) {
	for i := range tags {
		for j := range s.Items {
			if s.Items[j].Base.Id == tags[i].ResourceId {
				s.Items[j].Information.AddTag(tags[i])
			}
		}
	}
}

func (r *Information) AddTag(t *Tag) {
	r.Tags = append(r.Tags, t)
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

// NewSearchRequestFromHTTP keywords=xx&domain=xx&tag=app=~app1,app2,app3
func NewSearchRequestFromHTTP(r *http.Request) (*SearchRequest, error) {
	qs := r.URL.Query()
	req := &SearchRequest{
		Page:        NewPageRequestFromHTTP(r),
		Keywords:    qs.Get("keywords"),
		ExactMatch:  qs.Get("exact_match") == "true",
		Domain:      qs.Get("domain"),
		Namespace:   qs.Get("namespace"),
		Env:         qs.Get("env"),
		Status:      qs.Get("status"),
		SyncAccount: qs.Get("sync_account"),
		WithTags:    qs.Get("with_tags") == "true",
		Tags:        []*TagSelector{},
	}

	umStr := qs.Get("usage_mode")
	if umStr != "" {
		mode, err := ParseUsageModeFromString(umStr)
		if err != nil {
			return nil, err
		}
		req.UsageMode = &mode
	}

	rtStr := qs.Get("resource_type")
	if rtStr != "" {
		rt, err := ParseTypeFromString(rtStr)
		if err != nil {
			return nil, err
		}
		req.Type = &rt
	}

	// 单独处理Tag参数 app~=app1,app2,app3 --> TagSelector ---> req
	tgStr := qs.Get("tag")
	if tgStr != "" {
		tg, err := NewTagsFromString(tgStr)
		if err != nil {
			return nil, err
		}
		req.AddTag(tg...)
	}

	return req, nil
}

func (req *SearchRequest) AddTag(t ...*TagSelector) {
	req.Tags = append(req.Tags, t...)
}

// key1=v1,v2,v3&key2=~v1,v2,v3
func NewTagsFromString(tagStr string) (tags []*TagSelector, err error) {
	if tagStr == "" {
		return
	}

	items := strings.Split(tagStr, "&")
	for _, v := range items {
		// key1=v1,v2,v3 --> TagSelector
		t, err := ParExpr(v)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	return
}

func ParExpr(str string) (*TagSelector, error) {
	var (
		op = ""
		kv = []string{}
	)

	// app=~v1,v2,v3
	if strings.Contains(str, Operator_LIKE_EQUAL) {
		op = "LIKE"
		kv = strings.Split(str, Operator_LIKE_EQUAL)
	} else if strings.Contains(str, Operator_NOT_LIKE_EQUAL) {
		op = "NOT LIKE"
		kv = strings.Split(str, Operator_NOT_LIKE_EQUAL)
	} else if strings.Contains(str, Operator_NOT_EQUAL) {
		op = "!="
		kv = strings.Split(str, Operator_NOT_EQUAL)
	} else if strings.Contains(str, Operator_EQUAL) {
		op = "="
		kv = strings.Split(str, Operator_EQUAL)
	} else {
		return nil, fmt.Errorf("no support operator [=, =~, !=, !~]")
	}

	if len(kv) != 2 {
		return nil, fmt.Errorf("key,value format error, requred key=value")
	}

	selector := &TagSelector{
		Key:      kv[0],
		Operator: op,
		Values:   []string{},
	}

	// v1,v2,v3 splite [v1,v2,v3]
	// 如果Value等于*表示只匹配key
	if kv[1] != "*" {
		selector.Values = strings.Split(kv[1], ",")
	}

	return selector, nil
}

func (i *Information) Hash() string {
	return utils.Hash(i)
}

func (i *Information) PrivateIPToString() string {
	return strings.Join(i.PrivateIp, ",")
}

func (i *Information) PublicIPToString() string {
	return strings.Join(i.PublicIp, ",")
}

func NewThirdTag(key, value string) *Tag {
	return &Tag{
		Type:   TagType_THIRD,
		Key:    key,
		Value:  value,
		Weight: 1,
	}
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
