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

	project, err := handlers.DB.GetProject(projectUpdate.ProjectHandle, sourceProvider.Key())
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
		builds[i] = &models.Build{Handle: buildID}
	}

	hasErrs := false
	for _, build := range builds {
		if err = handlers.SendStartDeployMessage(build.Handle, sourceProvider.Key()); err != nil {
			logging.InternalError(err)
			hasErrs = true
		}
	}
	if hasErrs {
		r.RespondWithServerError(fmt.Errorf("Builds "))
		return
	}

	r.RespondWithSuccess()
}
