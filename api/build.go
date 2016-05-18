package api

import "github.com/mikec/msplapi/models"

type BuildReq struct {
	Handle        string  `json:"handle"`
	ProjectHandle string  `json:"project_handle"`
	Provider      string  `json:"provider"`
	ProviderUrl   *string `json:"provider_url,omitempty"`
}

func (br *BuildReq) SourceProvider() string {
	return br.Provider
}

func (handlers *Handlers) SaveBuild(r *Responder, ctx *RequestContext) {
	buildReq := ctx.RequestData.(*BuildReq)

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
