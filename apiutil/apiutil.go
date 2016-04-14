package apiutil

import (
	"github.com/mikec/marsupi-api/models"

  "io"
	"net/http"
	"fmt"
	"strings"
	"encoding/json"
)

type ApiUtil struct {
  ServerURL			string
}

func (self *ApiUtil) ProjectsURL() string {
	return fmt.Sprintf("%s/api/0/projects", self.ServerURL)
}

func (self *ApiUtil) ProjectURL(projectId int64) string {
	return fmt.Sprintf("%s/%d", self.ProjectsURL(), projectId)
}

func (self *ApiUtil) CreateProject(JSON string) (*http.Response, error) {
  res, err := DoPost(self.ProjectsURL(), JSON)
  if err != nil {
    return nil, err
  }
  return res, nil
}

func (self *ApiUtil) DeleteProject(id int64) (*http.Response, error) {
  res, err := DoDelete(self.ProjectURL(id))
  if err != nil {
    return nil, err
  }
  return res, nil
}

func (self *ApiUtil) GetProjects() ([]models.Project, *http.Response, error) {
  res, err := DoGet(self.ProjectsURL())
  if err != nil {
    return nil, nil, err
  }

  decoder := json.NewDecoder(res.Body)
  var projects []models.Project
  err = decoder.Decode(&projects)
  if err != nil {
    return nil, res, err
  }

  return projects, res, nil
}

func (self *ApiUtil) GetProject(id int64) (*models.Project, *http.Response, error) {
  res, err := DoGet(self.ProjectURL(id))
  if err != nil {
    return nil, nil, err
  }

  decoder := json.NewDecoder(res.Body)
  var p *models.Project
  err = decoder.Decode(&p)
  if err != nil {
    return nil, res, err
  }

  return p, res, nil
}

func DoGet(url string) (*http.Response, error) {
  return DoReq("GET", url, nil)
}

func DoPost(url string, JSON string) (*http.Response, error) {
  return DoReq("POST", url, &JSON)
}

func DoDelete(url string) (*http.Response, error) {
  return DoReq("DELETE", url, nil)
}

func DoReq(method string, url string, JSON *string) (*http.Response, error) {
  var reader io.Reader
  if JSON != nil {
    reader = strings.NewReader(*JSON)
  }

  req, err := http.NewRequest(method, url, reader)
  if err != nil {
    return nil, err
  }

  res, err := http.DefaultClient.Do(req)
  if err != nil {
    return nil, err
  }

  return res, nil
}

