package api

import (
    "net/http"
    "log"
    "fmt"
    "strconv"

    "github.com/gorilla/mux"

    "github.com/mikec/marsupi-api/models"
)

// curl -X POST http://localhost:8080/api/0/projects -d '{"service":"github", "username":"mikec", "name":"example-project"}'
func (handlers *Handlers) ProjectPost(w http.ResponseWriter, req *http.Request) {
    r := Responder{w}
    
    var p models.Project
    if err := DecodeRequest(req, &p); err != nil {
        log.Println(err)
        r.RespondWithError(InvalidJsonApiErr)
        return
    }
    
    if err := handlers.db.SaveProject(&p); err != nil {
        qErr := err.(*models.QueryExecError)
        switch qErr.Name {
        case "unique_violation":
            r.RespondWithError("Project already exists")
        default:
            r.RespondWithError("Failed to save project")
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
        log.Println(err)
        r.RespondWithError("Failed to get projects")
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
        apiErr := &ApiError{}
        apiErr.InvalidParam("ID", project_id)
        r.RespondWithError(*apiErr)
        return
      }

    project, err := handlers.db.GetProjectById(id)
    if err != nil {
        log.Println(err)
        r.RespondWithError(fmt.Sprintf("Failed to get project where id=%d", id))
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
        log.Println(err)
        r.RespondWithError(fmt.Sprintf("%s is not a valid project ID", project_id))
        return
    }

    if err := handlers.db.DeleteProject(id); err != nil {
        log.Println(err)
        r.RespondWithError("Failed to delete project")
        return
    }

    r.RespondWithSuccess()
}
