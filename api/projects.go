package api

import (
	"github.com/mikec/msplapi/models"
	"github.com/mikec/msplapi/provider"
)

type PostProjectReq struct {
	Handle   string `json:"handle" validate:"nonzero"`
	Provider string `json:"provider" validate:"nonzero"`
}

func (pr *PostProjectReq) SourceProvider() string {
	return pr.Provider
}

type ProjectResp struct {
	Handle   string `json:"handle"`
	Url      string `json:"url"`
	Provider string `json:"provider"`
}

func (handlers *Handlers) PostProject(r *Responder, ctx *RequestContext) {

	reqData := ctx.RequestData.(*PostProjectReq)

	sourceProvider := *ctx.SourceProvider
	projectData, err := sourceProvider.GetProjectData(
		ctx.SourceProviderToken.AccessToken,
		reqData.Handle,
	)
	if err != nil {
		switch v := err.(type) {
		case *provider.ProviderError:
			r.RespondWithError(NewApiErr(v.Error()))
		default:
			r.RespondWithServerError(err)
		}
		return
	}

	project := models.Project{
		Handle:         reqData.Handle,
		SourceUrl:      projectData.Url,
		SourceProvider: reqData.Provider,
	}

	err = handlers.db.SaveProject(&project, ctx.CurrentUser)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}

	r.Respond(&ProjectResp{
		project.Handle,
		project.SourceUrl,
		project.SourceProvider,
	})
}

func (handlers *Handlers) GetProviderProjects(r *Responder, ctx *RequestContext) {
	sourceProvider := *ctx.SourceProvider
	projects, err := sourceProvider.GetProjects(
		ctx.SourceProviderToken.AccessToken,
		ctx.SourceProviderToken.ProviderUsername,
	)
	if err != nil {
		switch v := err.(type) {
		case *provider.ProviderError:
			r.RespondWithError(NewApiErr(v.Error()))
		default:
			r.RespondWithServerError(err)
		}
		return
	}
	r.Respond(projects)
}

func (handlers *Handlers) GetProjects(r *Responder, ctx *RequestContext) {
	projects, err := handlers.db.GetUserProjects(ctx.CurrentUser)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}
	r.Respond(projects)
}
