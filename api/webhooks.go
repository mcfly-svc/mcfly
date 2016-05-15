package api

import (
	"fmt"

	"github.com/mikec/msplapi/logging"
	"github.com/mikec/msplapi/models"
)

func (handlers *Handlers) ProjectUpdateWebhook(r *Responder, ctx *RequestContext) {
	sourceProvider := *ctx.SourceProvider

	projectUpdate, err := sourceProvider.DecodeProjectUpdateRequest(r.Request)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}

	project, err := handlers.db.GetProject(projectUpdate.ProjectHandle, sourceProvider.Key())
	if err != nil {
		r.RespondWithServerError(err)
		return
	}
	if project == nil {
		err = fmt.Errorf(NewNotFoundErr("project", projectUpdate.ProjectHandle).Error)
		r.RespondWithServerError(err)
		return
	}

	builds := make([]*models.Build, len(projectUpdate.Builds))
	for i, buildID := range projectUpdate.Builds {
		builds[i] = &models.Build{
			ProjectID: project.ID,
			Handle:    buildID,
		}
	}

	// TODO NEXT: figure out what how to actually get the project code and deploy it
	// and add build commands to mclovin

	if errs := handlers.db.SaveBuilds(builds); errs != nil {
		for _, e := range errs {
			logging.InternalError(e)
		}
		r.RespondWithServerError(fmt.Errorf("SaveBuilds failed in ProjectUpdateWebhook"))
		return
	}

	r.RespondWithSuccess()
}
