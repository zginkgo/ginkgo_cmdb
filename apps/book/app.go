package book

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/imdario/mergo"
	"github.com/rs/xid"
)

const (
	AppName = "book"
)

const (
	// DefaultPageSize 默认分页大小
	DefaultPageSize = 20
	// DefaultPageNumber 默认页号
	DefaultPageNumber = 1
)

var (
	validate = validator.New()
)

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

func NewCreateBookRequest() *CreateBookRequest {
	return &CreateBookRequest{}
}

func NewBook(req *CreateBookRequest) (*Book, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &Book{
		Id:       xid.New().String(),
		CreateAt: time.Now().UnixMicro(),
		Data:     req,
	}, nil
}

func (req *CreateBookRequest) Validate() error {
	return validate.Struct(req)
}

func NewBookSet() *BookSet {
	return &BookSet{
		Items: []*Book{},
	}
}

func (s *BookSet) Add(item *Book) {
	s.Items = append(s.Items, item)
}

func NewDefaultBook() *Book {
	return &Book{
		Data: &CreateBookRequest{},
	}
}

func (i *Book) Update(req *UpdateBookRequest) {
	i.UpdateAt = time.Now().UnixMicro()
	i.UpdateBy = req.UpdateBy
	i.Data = req.Data
}

func (i *Book) Patch(req *UpdateBookRequest) error {
	i.UpdateAt = time.Now().UnixMicro()
	i.UpdateBy = req.UpdateBy
	return mergo.MergeWithOverwrite(i.Data, req.Data)
}

func NewDescribeBookRequest(id string) *DescribeBookRequest {
	return &DescribeBookRequest{
		Id: id,
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

func NewQueryBookRequest() *QueryBookRequest {
	return &QueryBookRequest{
		Page: NewDefaultPageRequest(),
	}
}

func NewQueryBookRequestFromHTTP(r *http.Request) *QueryBookRequest {
	qs := r.URL.Query()

	return &QueryBookRequest{
		Page:     NewPageRequestFromHTTP(r),
		Keywords: qs.Get("keywords"),
	}
}

func NewPutBookRequest(id string) *UpdateBookRequest {
	return &UpdateBookRequest{
		Id:         id,
		UpdateMode: UpdateMode_PUT,
		UpdateAt:   time.Now().UnixMicro(),
		Data:       NewCreateBookRequest(),
	}
}

func NewPatchBookRequest(id string) *UpdateBookRequest {
	return &UpdateBookRequest{
		Id:         id,
		UpdateMode: UpdateMode_PATCH,
		UpdateAt:   time.Now().UnixMicro(),
		Data:       NewCreateBookRequest(),
	}
}

func NewDeleteBookRequestWithID(id string) *DeleteBookRequest {
	return &DeleteBookRequest{
		Id: id,
	}
}
