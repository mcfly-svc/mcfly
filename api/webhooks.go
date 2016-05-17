package api

import (
	"fmt"

	"github.com/mikec/msplapi/logging"
	"github.com/mikec/msplapi/models"
	"github.com/mikec/msplapi/mq"
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
		dpq := mq.DeployQueueMessage{
			BuildHandle:               build.Handle,
			ProjectHandle:             project.Handle,
			SourceProvider:            sourceProvider.Key(),
			SourceProviderAccessToken: "", // TODO: get this after verifying the webhook request
		}
		if err = handlers.MessageChannel.SendDeployQueueMessage(&dpq); err != nil {
			logging.InternalError(err)
			hasErrs = true
		}
	}
	if hasErrs {
		r.RespondWithServerError(fmt.Errorf("Calls to SendDeployQueueMessage failed"))
		return
	}

	r.RespondWithSuccess()
}
