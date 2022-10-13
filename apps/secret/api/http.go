package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/zginkgo/ginkgo_cmdb/apps"
	"github.com/zginkgo/ginkgo_cmdb/apps/secret"
)

var (
	h = &handler{}
)

type handler struct {
	service secret.ServiceServer
	log     logger.Logger
}

//func (h *handler) Config() error {
//	//h.log = zap.L().Named("Secret")
//	//if apps.Secret == nil {
//	//	return fmt.Errorf("dependence service secret not ready")
//	//}
//	//h.service = apps.Secret
//	h.log = zap.L().Named(secret.AppName)
//	h.service = apps.GetHttpApp(secret.AppName).(secret.ServiceServer)
//	return nil
//}

func (h *handler) Config() error {
	h.log = zap.L().Named(secret.AppName)
	h.service = apps.GetGrpcApp(secret.AppName).(secret.ServiceServer)
	return nil
}

func (h *handler) Name() string {
	return secret.AppName
}

func (h *handler) Version() string {
	return "v1"
}

//func (h *handler) Registry(r router.SubRouter) {
//	r.Handle("POST", "/secrets", h.CreateSecret)
//	r.Handle("GET", "/secrets", h.QuerySecret)
//	r.Handle("GET", "/secrets/:id", h.DescribeSecret)
//	r.Handle("DELETE", "/secrets/:id", h.DeleteSecret)
//}

func (h *handler) Registry(ws *restful.WebService) {
	tags := []string{h.Name()}

	// Route /cmdb/api/v1/book/ --> h.CreateBook
	// 需要装饰路由:   Route Path 不需要认证
	ws.Route(ws.POST("/").To(h.CreateSecret).
		Doc("create a secret").
		Metadata(restfulspec.KeyOpenAPITags, tags).

		// 添加一下标签做 路由装饰
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Create.Value()).

		// 是否开启认证
		Metadata(label.Auth, label.Enable).

		// 是否开启鉴权
		Metadata(label.Permission, label.Enable).

		// 基于用户属性的权限装饰
		Metadata(label.Allow, "admin").
		Reads(secret.CreateSecretRequest{}).
		Writes(response.NewData(secret.Secret{})))

	ws.Route(ws.GET("/").To(h.QuerySecret).
		Doc("get all secret").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.List.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Reads(secret.QuerySecretRequest{}).
		Writes(response.NewData(secret.SecretSet{})).
		Returns(200, "OK", secret.SecretSet{}))

	ws.Route(ws.GET("/{id}").To(h.DescribeSecret).
		Doc("describe an secret").
		Param(ws.PathParameter("id", "identifier of the secret").DataType("integer").DefaultValue("1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Get.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Writes(response.NewData(secret.Secret{})).
		Returns(200, "OK", response.NewData(secret.Secret{})).
		Returns(404, "Not Found", nil))

	ws.Route(ws.DELETE("/{id}").To(h.DeleteSecret).
		Doc("delete a secret").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata(label.Resource, h.Name()).
		Metadata(label.Action, label.Delete.Value()).
		Metadata(label.Auth, label.Enable).
		Metadata(label.Permission, label.Enable).
		Param(ws.PathParameter("id", "identifier of the secret").DataType("string")))
}

func init() {
	apps.RegistryRESTfulApp(h)
}
