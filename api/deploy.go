package api

import "github.com/mikec/msplapi/mq"

type DeployReq struct {
	BuildHandle         string `json:"build_handle" validate:"nonzero"`
	SourceProjectHandle string `json:"project_handle" validate:"nonzero"`
	Provider            string `json:"provider" validate:"nonzero"`
}

func (r *DeployReq) SourceProvider() string {
	return r.Provider
}

func (r *DeployReq) ProjectHandle() string {
	return r.SourceProjectHandle
}

func (handlers *Handlers) Deploy(r *Responder, ctx *RequestContext) {
	deployReq := ctx.RequestData.(*DeployReq)
	sp := *ctx.SourceProvider

	dpq := mq.DeployQueueMessage{
		BuildHandle:               deployReq.BuildHandle,
		ProjectHandle:             ctx.Project.Handle,
		SourceProvider:            sp.Key(),
		SourceProviderAccessToken: ctx.SourceProviderToken.AccessToken,
	}
	err := handlers.MessageChannel.SendDeployQueueMessage(&dpq)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}

	r.RespondWithSuccess()
}
