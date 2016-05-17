package api

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

	err := handlers.SendStartDeployMessage(deployReq.BuildHandle, sp.Key())
	if err != nil {
		r.RespondWithServerError(err)
		return
	}

	r.RespondWithSuccess()
}

func (handlers *Handlers) SendStartDeployMessage(buildHandle, provider string) error {
	return handlers.MessageChannel.StartDeploy(buildHandle, provider)
}
