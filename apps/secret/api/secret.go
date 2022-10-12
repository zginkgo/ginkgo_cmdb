package api

import (
	"github.com/zginkgo/ginkgo_cmdb/apps/secret"
	"github.com/zginkgo/ginkgo_cmdb/utils"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/response"
)

func (h *handler) CreateSecret(r *restful.Request, w *restful.Response) {
	req := secret.NewCreateSecretRequest()
	if err := utils.GetDataFromRequest(r.Request, req); err != nil {
		response.Failed(w, err)
		return
	}

	set, err := h.service.CreateSecret(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	response.Success(w.ResponseWriter, set)
}

func (h *handler) QuerySecret(r *restful.Request, w *restful.Response) {
	req := secret.NewQuerySecretRequestFromHTTP(r.Request)
	set, err := h.service.QuerySecret(r.Request.Context(), req)

	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}

func (h *handler) DescribeSecret(r *restful.Request, w *restful.Response) {
	req := secret.NewDescribeSecretRequest(r.PathParameter("id"))
	ins, err := h.service.DescribeSecret(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	ins.Data.Desense()
	response.Success(w.ResponseWriter, ins)
}

func (h *handler) DeleteSecret(r *restful.Request, w *restful.Response) {
	req := secret.NewDeleteSecretRequestWithID(r.PathParameter("id"))
	set, err := h.service.DeleteSecret(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}
