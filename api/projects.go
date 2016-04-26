package api

import (
    "net/http"
    "strconv"

    "github.com/gorilla/mux"

    "github.com/mikec/marsupi-api/models"
    "github.com/mikec/marsupi-api/logging"
)

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
