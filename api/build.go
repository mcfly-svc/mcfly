package api

import (
	"github.com/mcfly-svc/mcfly/api/apidata"
	"github.com/mcfly-svc/mcfly/models"
)

func (handlers *Handlers) SaveBuild(r *Responder, ctx *RequestContext) {
	buildReq := ctx.RequestData.(*apidata.BuildReq)

	p, err := handlers.DB.GetProject(buildReq.ProjectHandle, buildReq.Provider)
	if err != nil {
		r.RespondWithError(NewNotFoundErr("project", buildReq.ProjectHandle))
		return
	}

	build := &models.Build{Handle: buildReq.Handle, ProviderUrl: buildReq.ProviderUrl}

	if err := handlers.DB.SaveBuild(build, p); err != nil {
		r.RespondWithServerError(err)
		return
	}

	r.RespondWithSuccess()
}
