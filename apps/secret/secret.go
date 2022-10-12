package secret

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zginkgo/ginkgo_cmdb/conf"
	"github.com/go-playground/validator/v10"
	"github.com/infraboard/mcube/crypto/cbc"
	"github.com/rs/xid"
)

const (
	AppName = "secret"
	// DefaultPageSize 默认分页大小
	DefaultPageSize = 20
	// DefaultPageNumber 默认页号
	DefaultPageNumber = 1
)

var (
	validate = validator.New()
)

func NewDefaultSecret() *Secret {
	return &Secret{
		Data: &CreateSecretRequest{
			RequestRate: 5,
		},
	}
}

func NewSecretSet() *SecretSet {
	return &SecretSet{
		Items: []*Secret{},
	}
}

func (s *SecretSet) Add(item *Secret) {
	s.Items = append(s.Items, item)
}

func NewSecret(req *CreateSecretRequest) (*Secret, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &Secret{
		Id:       xid.New().String(),
		CreateAt: time.Now().UnixMilli(),
		Data:     req,
	}, nil
}

func NewCreateSecretRequest() *CreateSecretRequest {
	return &CreateSecretRequest{
		RequestRate: 5,
	}
}

func (req *CreateSecretRequest) Validate() error {
	if len(req.AllowRegions) == 0 {
		return fmt.Errorf("required less one allow_regions")
	}
	return validate.Struct(req)
}

func NewQuerySecretRequestFromHTTP(r *http.Request) *QuerySecretRequest {
	qs := r.URL.Query()

	return &QuerySecretRequest{
		Page:     NewPageRequestFromHTTP(r),
		Keywords: qs.Get("keywords"),
	}
}

func NewQuerySecretRequest() *QuerySecretRequest {
	return &QuerySecretRequest{
		Page:     NewDefaultPageRequest(),
		Keywords: "",
	}
}

func NewDescribeSecretRequest(id string) *DescribeSecretRequest {
	return &DescribeSecretRequest{
		Id: id,
	}
}

func NewDeleteSecretRequestWithID(id string) *DeleteSecretRequest {
	return &DeleteSecretRequest{
		Id: id,
	}
}

func (s *CreateSecretRequest) AllowRegionString() string {
	return strings.Join(s.AllowRegions, ",")
}

func (s *CreateSecretRequest) LoadAllowRegionFromString(regions string) {
	if regions != "" {
		s.AllowRegions = strings.Split(regions, ",")
	}
}

func (s *CreateSecretRequest) EncryptAPISecret(key string) error {
	// 判断文本是否已经加密
	if strings.HasPrefix(s.ApiSecret, conf.C().App.EncryptKey) {
		return fmt.Errorf("text has ciphered")
	}

	cipherText, err := cbc.Encrypt([]byte(s.ApiSecret), []byte(key))
	if err != nil {
		return err
	}

	base64Str := base64.StdEncoding.EncodeToString(cipherText)
	s.ApiSecret = fmt.Sprintf("%s%s", conf.C().App.EncryptKey, base64Str)
	return nil
}

func (s *CreateSecretRequest) DecryptAPISecret(key string) error {
	// 判断文本是否已经是明文
	if !strings.HasPrefix(s.ApiSecret, conf.C().App.EncryptKey) {
		return fmt.Errorf("text is plan text")
	}

	base64CipherText := strings.TrimPrefix(s.ApiSecret, conf.C().App.EncryptKey)

	cipherText, err := base64.StdEncoding.DecodeString(base64CipherText)
	if err != nil {
		return err
	}

	planText, err := cbc.Decrypt([]byte(cipherText), []byte(key))
	if err != nil {
		return err
	}

	s.ApiSecret = string(planText)
	return nil
}

// 敏感数据脱敏
func (s *CreateSecretRequest) Desense() {
	if s.ApiSecret != "" {
		s.ApiSecret = "******"
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
