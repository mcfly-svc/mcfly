package api

import (
	"fmt"

	"github.com/mikec/msplapi/models"
)

func (handlers *Handlers) ProjectUpdateWebhook(r *Responder, ctx *RequestContext) {
	sourceProvider := *ctx.SourceProvider

	projectUpdate, err := sourceProvider.DecodeProjectUpdateRequest(r.Request)
	if err != nil {
		r.RespondWithServerError(err)
	}

	project, err := handlers.db.GetProject(projectUpdate.ProjectHandle, sourceProvider.Key())
	if err != nil {
		r.RespondWithServerError(err)
	}
	if project == nil {
		err = fmt.Errorf(NewNotFoundErr("project", projectUpdate.ProjectHandle).Error)
		r.RespondWithServerError(err)
	}

	builds := make([]*models.Build, len(projectUpdate.Builds))
	for i, buildID := range projectUpdate.Builds {
		builds[i] = &models.Build{
			ProjectID: project.ID,
			Handle:    buildID,
		}
	}

	if err = handlers.db.SaveBuilds(builds); err != nil {
		r.RespondWithServerError(err)
	}

	r.RespondWithSuccess()
}
