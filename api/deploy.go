package api

import (
	"github.com/mcfly-svc/mcfly/api/apidata"
	"github.com/mcfly-svc/mcfly/mq"
)

func (handlers *Handlers) Deploy(r *Responder, ctx *RequestContext) {
	deployReq := ctx.RequestData.(*apidata.DeployReq)
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
