package api

import (
    "net/http"
    "log"
    "fmt"
    "encoding/json"
    "strconv"

    "github.com/gorilla/mux"

    "github.com/mikec/marsupi-api/models"
)

// curl -X POST http://localhost:8080/api/0/projects -d '{"service":"github.com", "username":"mikec", "name":"example-project"}'
func (handlers *Handlers) ProjectsPost(w http.ResponseWriter, req *http.Request) {
    
    var p *models.Project
    if err := decodeRequestBodyJson(req, &p); err != nil {
        log.Println(err)
        writeErrorMessage(w, "Invalid JSON")
        return
    }
    
    if saveErr := handlers.db.SaveProject(p); saveErr != nil {
        log.Println(saveErr)
        writeErrorMessage(w, "Failed to save project")
        return
    }

    writeSuccessResponse(w)
}

// curl -X GET http://localhost:8080/api/0/projects
func (handlers *Handlers) ProjectsGet(w http.ResponseWriter, req *http.Request) {

    projects, err := handlers.db.GetProjects()
    if err != nil {
        log.Println(err)
        writeErrorMessage(w, "Failed to get projects")
        return
    }

    if err := json.NewEncoder(w).Encode(projects); err != nil {
        log.Fatal(err)
    }
}

// curl -X DELETE http://localhost:8080/api/0/projects/1
func (handlers *Handlers) ProjectsDelete(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    project_id := vars["project_id"]
    id, err := strconv.ParseInt(project_id, 10, 64)
    if err != nil {
        log.Println(err)
        writeErrorMessage(w, fmt.Sprintf("%s is not a valid project ID", project_id))
        return
    }

    if err := handlers.db.DeleteProject(id); err != nil {
        log.Println(err)
        writeErrorMessage(w, "Failed to delete project")
        return
    }

    writeSuccessResponse(w)
}
