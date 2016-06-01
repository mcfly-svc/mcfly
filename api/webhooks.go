package api

import (
	"fmt"

	"github.com/mcfly-svc/mcfly/logging"
	"github.com/mcfly-svc/mcfly/models"
	"github.com/mcfly-svc/mcfly/mq"
)

func (handlers *Handlers) ProjectUpdateWebhook(r *Responder, ctx *RequestContext) {
	sourceProvider := *ctx.SourceProvider
	spKey := sourceProvider.Key()

	projectUpdate, err := sourceProvider.DecodeProjectUpdateRequest(r.Request)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}

	project, err := handlers.DB.GetProject(projectUpdate.ProjectHandle, spKey)
	nfe := NewNotFoundErr("project", projectUpdate.ProjectHandle).Error
	if respondToNotFoundErr(r, nfe, err) {
		return
	}

	user, err := handlers.DB.GetUserByProject(project)
	nfe = fmt.Sprintf("GetUserByProject failed: could not find %s project %s", spKey, project.Handle)
	if respondToNotFoundErr(r, nfe, err) {
		return
	}

	providerAccessToken, err := handlers.DB.GetProviderTokenForUser(user, spKey)
	nfe = fmt.Sprintf("No %s access token for user %s", spKey, user.ID)
	if respondToNotFoundErr(r, nfe, err) {
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
			SourceProvider:            spKey,
			SourceProviderAccessToken: providerAccessToken.AccessToken,
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

func respondToNotFoundErr(r *Responder, notFoundMsg string, err error) bool {
	if err != nil {
		if err == models.ErrNotFound {
			r.RespondWithServerError(fmt.Errorf(notFoundMsg))
			return true
		}
		r.RespondWithServerError(err)
		return true
	}
	return false
}
