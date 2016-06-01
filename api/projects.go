package api

import (
	"github.com/mcfly-svc/mcfly/api/apidata"
	"github.com/mcfly-svc/mcfly/logging"
	"github.com/mcfly-svc/mcfly/models"
	"github.com/mcfly-svc/mcfly/provider"
)

func (handlers *Handlers) PostProject(r *Responder, ctx *RequestContext) {

	reqData := ctx.RequestData.(*apidata.ProjectReq)

	sourceProvider := *ctx.SourceProvider
	token := ctx.SourceProviderToken.AccessToken
	projectData, err := sourceProvider.GetProjectData(token, reqData.Handle)
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

	err = handlers.DB.SaveProject(&project, ctx.CurrentUser)
	if err != nil {
		if err == models.ErrDuplicate {
			r.RespondWithError(NewDuplicateErr("project", reqData.Handle))
		} else {
			r.RespondWithServerError(err)
		}
		return
	}

	r.Respond(&apidata.ProjectResp{
		project.Handle,
		project.SourceUrl,
		project.SourceProvider,
	})

	if err = sourceProvider.CreateProjectUpdateHook(token, reqData.Handle); err != nil {
		logging.InternalError(err)
	}
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
			if v.Code == "get_projects_failed" {
				// TODO: find out why this failed. If it's because the user's
				// github account was deleted, handle that and respond accordingly
				logging.InternalError(v)
			}
			r.RespondWithError(NewApiErr(v.Error()))
		default:
			r.RespondWithServerError(err)
		}
		return
	}
	r.Respond(projects)
}

func (handlers *Handlers) GetProjects(r *Responder, ctx *RequestContext) {
	projects, err := handlers.DB.GetUserProjects(ctx.CurrentUser)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}
	projectsResp := make([]apidata.ProjectResp, len(projects))
	for i, p := range projects {
		projectsResp[i] = apidata.ProjectResp{p.Handle, p.SourceUrl, p.SourceProvider}
	}
	r.Respond(projectsResp)
}

func (handlers *Handlers) DeleteProject(r *Responder, ctx *RequestContext) {
	project := ctx.RequestData.(*apidata.ProjectReq)

	err := handlers.DB.DeleteUserProject(ctx.CurrentUser, project.Provider, project.Handle)
	if err != nil {
		if err == models.ErrNotFound {
			r.RespondWithError(NewNotFoundErr("project", project.Handle))
		} else {
			r.RespondWithServerError(err)
		}
		return
	}

	r.RespondWithSuccess()
}
