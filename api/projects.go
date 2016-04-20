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
func (handlers *Handlers) ProjectsPost(w http.ResponseWriter, req *http.Request) {
    
    var p *models.Project
    if err := DecodeRequest(req, &p); err != nil {
        log.Println(err)
        writeErrorResponse(w, InvalidJsonApiErr)
        return
    }
    
    if saveErr := handlers.db.SaveProject(p); saveErr != nil {
        log.Println(saveErr)
        writeErrorResponse(w, "Failed to save projects")
        return
    }

    writeSuccessResponse(w)
}

// curl -X GET http://localhost:8080/api/0/projects
func (handlers *Handlers) ProjectsGet(w http.ResponseWriter, req *http.Request) {

    projects, err := handlers.db.GetProjects()
    if err != nil {
        log.Println(err)
        writeErrorResponse(w, "Failed to get projects")
        return
    }

    writeResponse(w, projects)
}

// curl -X GET http://localhost:8080/api/0/projects/1
func (handlers *Handlers) ProjectGet(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    project_id := vars["project_id"]
    id, err := strconv.ParseInt(project_id, 10, 64)

    project, err := handlers.db.GetProjectById(id)
    if err != nil {
        log.Println(err)
        writeErrorResponse(w, fmt.Sprintf("Failed to get project with id=%d", id))
        return
    }

    writeResponse(w, project)

}

// curl -X DELETE http://localhost:8080/api/0/projects/1
func (handlers *Handlers) ProjectsDelete(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    project_id := vars["project_id"]
    id, err := strconv.ParseInt(project_id, 10, 64)
    if err != nil {
        log.Println(err)
        writeErrorResponse(w, fmt.Sprintf("%s is not a valid project ID", project_id))
        return
    }

    if err := handlers.db.DeleteProject(id); err != nil {
        log.Println(err)
        writeErrorResponse(w, "Failed to delete project")
        return
    }

    writeSuccessResponse(w)
}

type ProjectDecoder struct {}

func (self ProjectDecoder) DecodeRequest(req *http.Request) (interface{}, error) {
    var p models.Project
    if err := DecodeRequest(req, &p); err != nil {
        return nil, err
    }
    return p, nil
}

func (self ProjectDecoder) DecodeResponse(res *http.Response) (interface{}, error) {
    var p models.Project
    if err := DecodeResponse(res, &p); err != nil {
        return nil, err
    }
    return p, nil
}

func (self ProjectDecoder) DecodeArrayRequest(req *http.Request) ([]interface{}, error) {
    var projects []models.Project
    if err := DecodeRequest(req, &projects); err != nil {
        return nil, err
    }
    ret := make([]interface{}, len(projects))
    for i, d := range projects {
        ret[i] = d
    }
    return ret, nil
}

func (self ProjectDecoder) DecodeArrayResponse(res *http.Response) ([]interface{}, error) {
    var projects []models.Project
    if err := DecodeResponse(res, &projects); err != nil {
        return nil, err
    }
    ret := make([]interface{}, len(projects))
    for i, d := range projects {
        ret[i] = d
    }
    return ret, nil
}
