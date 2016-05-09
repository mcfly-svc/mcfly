package api

import (
	"net/http"

	"github.com/mikec/msplapi/models"
	"github.com/mikec/msplapi/provider"
)

type PostProjectReq struct {
	ProjectHandle string `json:"project_handle" validate:"nonzero"`
	Provider      string `json:"provider" validate:"nonzero"`
}

func (handlers *Handlers) PostProject(w http.ResponseWriter, req *http.Request) {
	r := &Responder{w, req}

	user := r.ValidateAuthorization(handlers.db)
	if user == nil {
		return
	}

	var reqData PostProjectReq
	decodeErr := r.DecodeRequest(&reqData)
	if decodeErr != nil {
		return
	}

	reqValid := r.ValidateRequestData(&reqData)
	if !reqValid {
		return
	}

	sourceProvider := handlers.sourceProviders[reqData.Provider]
	if sourceProvider == nil {
		// TODO: change these errors to provider type specific
		r.RespondWithError(NewUnsupportedProviderErr(reqData.Provider))
		return
	}

	providerToken, err := handlers.db.GetProviderTokenForUser(user, reqData.Provider)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}
	if providerToken == nil {
		r.RespondWithError(NewProviderTokenNotFoundErr(reqData.Provider))
		return
	}

	projectData, err := sourceProvider.GetProjectData(*providerToken, reqData.ProjectHandle)
	if err != nil {
		pErr, ok := err.(*provider.GetProjectDataError)
		if !ok {
			r.RespondWithServerError(err)
			return
		}
		r.RespondWithError(NewApiErr(pErr.Error()))
		return
	}

	project := models.Project{
		Handle:         reqData.ProjectHandle,
		SourceUrl:      projectData.Url,
		SourceProvider: reqData.Provider,
	}
	err = handlers.db.SaveProject(&project, user)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}

	r.RespondWithSuccess()
}

/*
// curl -X POST http://localhost:8080/api/0/projects -d '{"service":"github", "username":"mikec", "name":"example-project"}'
func (handlers *Handlers) ProjectPost(w http.ResponseWriter, req *http.Request) {
	r := Responder{w}

	var p models.Project
	if err := DecodeRequest(req, &p); err != nil {
		logging.LogInternalError("ProjectPost", err)
		r.RespondWithError(NewInvalidJsonErr())
		return
	}

	if err := handlers.db.SaveProject(&p); err != nil {
		qErr := err.(*models.QueryExecError)
		switch qErr.Name {
		case "unique_violation":
			r.RespondWithError(NewDuplicateCreateErr("project"))
		default:
			logging.LogInternalError("ProjectPost", err)
			r.RespondWithError(NewCreateFailedErr("project"))
		}
		return
	}

	r.Respond(p)
}

// curl -X GET http://localhost:8080/api/0/projects
func (handlers *Handlers) ProjectsGet(w http.ResponseWriter, req *http.Request) {
	r := Responder{w}

	projects, err := handlers.db.GetProjects()
	if err != nil {
		logging.LogInternalError("ProjectsGet", err)
		r.RespondWithError(NewGetEntitiesErr("projects"))
		return
	}

	r.Respond(projects)
}

// curl -X GET http://localhost:8080/api/0/projects/1
func (handlers *Handlers) ProjectGet(w http.ResponseWriter, req *http.Request) {
	r := Responder{w}

	vars := mux.Vars(req)
	project_id := vars["project_id"]
	id, err := strconv.ParseInt(project_id, 10, 64)

	if id == 0 || err != nil {
		r.RespondWithError(NewInvalidParamErr("ID", project_id))
		return
	}

	project, err := handlers.db.GetProjectById(id)
	if err != nil {
		logging.LogInternalError("ProjectGet", err)
		r.RespondWithError(NewGetEntityErr("project", id))
		return
	}

	r.Respond(project)
}

// curl -X DELETE http://localhost:8080/api/0/projects/1
func (handlers *Handlers) ProjectDelete(w http.ResponseWriter, req *http.Request) {
	r := Responder{w}

	vars := mux.Vars(req)
	project_id := vars["project_id"]
	id, err := strconv.ParseInt(project_id, 10, 64)

	if err != nil {
		logging.LogInternalError("ProjectDelete", err)
		r.RespondWithError(NewInvalidParamErr("project_id", project_id))
		return
	}

	if err := handlers.db.DeleteProject(id); err != nil {
		logging.LogInternalError("ProjectDelete", err)
		r.RespondWithError(NewDeleteFailedErr("project"))
		return
	}

	r.RespondWithSuccess()
}
*/
