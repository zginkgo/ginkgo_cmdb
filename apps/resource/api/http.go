package api

import (
	"github.com/zginkgo/ginkgo_cmdb/apps"
	"github.com/zginkgo/ginkgo_cmdb/apps/resource"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	h = &handler{}
)

type handler struct {
	service resource.ServiceServer
	log     logger.Logger
}

func (h *handler) Config() error {
	h.log = zap.L().Named(resource.AppName)
	h.service = apps.GetGrpcApp(resource.AppName).(resource.ServiceServer)
	return nil
}

func (h *handler) Name() string {
	return resource.AppName
}

func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(ws *restful.WebService) {
	ws.Route(ws.GET("/search").To(h.SearchResource))

}

func init() {
	apps.RegistryRESTfulApp(h)
}
