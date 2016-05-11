package api

import (
	"github.com/mikec/msplapi/models"
	"github.com/mikec/msplapi/provider"
)

type ProjectReq struct {
	Handle   string `json:"handle" validate:"nonzero"`
	Provider string `json:"provider" validate:"nonzero"`
}

func (pr *ProjectReq) SourceProvider() string {
	return pr.Provider
}

type ProjectResp struct {
	Handle   string `json:"handle"`
	Url      string `json:"url"`
	Provider string `json:"provider"`
}

func (handlers *Handlers) PostProject(r *Responder, ctx *RequestContext) {

	reqData := ctx.RequestData.(*ProjectReq)

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
	projectsResp := make([]ProjectResp, len(projects))
	for i, p := range projects {
		projectsResp[i] = ProjectResp{p.Handle, p.SourceUrl, p.SourceProvider}
	}
	r.Respond(projectsResp)
}

func (handlers *Handlers) DeleteProject(r *Responder, ctx *RequestContext) {
	project := ctx.RequestData.(*ProjectReq)

	err := handlers.db.DeleteUserProject(ctx.CurrentUser, project.Provider, project.Handle)
	if err != nil {
		switch v := err.(type) {
		case *models.ModelsError:
			if v.Code == "not_found" {
				r.RespondWithError(NewNotFoundErr("project", project.Handle))
			} else {
				r.RespondWithServerError(err)
			}
		default:
			r.RespondWithServerError(err)
		}
		return
	}

	r.RespondWithSuccess()
}
