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

	var reqData PostProjectReq
	reqData = *ctx.RequestData.(*PostProjectReq)

	providerToken, err := handlers.db.GetProviderTokenForUser(ctx.CurrentUser, reqData.Provider)
	if err != nil {
		r.RespondWithServerError(err)
		return
	}
	if providerToken == nil {
		r.RespondWithError(NewProviderTokenNotFoundErr(reqData.Provider))
		return
	}

	sourceProvider := *ctx.SourceProvider
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
