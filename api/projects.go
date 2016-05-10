package api

import (
	"github.com/mikec/msplapi/models"
	"github.com/mikec/msplapi/provider"
)

type PostProjectReq struct {
	ProjectHandle string `json:"project_handle" validate:"nonzero"`
	Provider      string `json:"provider" validate:"nonzero"`
}

func (pr *PostProjectReq) SourceProvider() string {
	return pr.Provider
}

type PostProjectResp struct {
	ProjectHandle string `json:"project_handle"`
	ProjectUrl    string `json:"project_url"`
	Provider      string `json:"provider"`
}

func (handlers *Handlers) PostProject(r *Responder, ctx *RequestContext) {

	reqData := ctx.RequestData.(*PostProjectReq)

	sourceProvider := *ctx.SourceProvider
	projectData, err := sourceProvider.GetProjectData(
		ctx.SourceProviderToken.AccessToken,
		reqData.ProjectHandle,
	)
	if err != nil {
		switch v := err.(type) {
		case *provider.GetProjectDataError:
			r.RespondWithError(NewApiErr(v.Error()))
		case *provider.ProviderTokenInvalidError:
			r.RespondWithError(NewApiErr(v.Error()))
		default:
			r.RespondWithServerError(err)
		}
		return
	}

	project := models.Project{
		Handle:         reqData.ProjectHandle,
		SourceUrl:      projectData.Url,
		SourceProvider: reqData.Provider,
	}

	err = handlers.db.SaveProject(&project, ctx.CurrentUser)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}

	r.Respond(&PostProjectResp{
		project.Handle,
		project.SourceUrl,
		project.SourceProvider,
	})
}

type GetProjectsResp []struct {
	ProjectHandle string `json:"project_handle"`
	ProjectUrl    string `json:"project_url"`
}

func (handlers *Handlers) GetProjects(r *Responder, ctx *RequestContext) {
	sourceProvider := *ctx.SourceProvider
	projects, err := sourceProvider.GetProjects(
		ctx.SourceProviderToken.AccessToken,
		ctx.SourceProviderToken.ProviderUsername,
	)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}
	r.Respond(projects)
}
