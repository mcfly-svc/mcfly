package api_test

import (
	"github.com/mikec/marsupi-api/models"

	"testing"
  "fmt"
  "encoding/json"
  "net/http"
  "strings"
  "io/ioutil"
)

// create project should return 200 status
func TestCreateProject(t *testing.T) {
	cleanupDB()

	res := createExampleProject(t)

  if res.StatusCode != 200 {
      t.Errorf("Success expected: Got %d", res.StatusCode)
  }
}

// create 3 projects, and expect to get 3 projects
func TestGetProjects(t *testing.T) {
	cleanupDB()

	createProject(t, `{"service": "github", "username": "mikec", "name": "example-project-1"}`)
	createProject(t, `{"service": "github", "username": "mikec", "name": "example-project-2"}`)
	createProject(t, `{"service": "github", "username": "mikec", "name": "example-project-3"}`)

	projects := getProjects(t)

  expected := 3
  actual := len(projects)
  if expected != actual {
  	t.Error(fmt.Sprintf("Expected %d projects but got %d", expected, actual))
  }

}

// creating two projecs with the same service/username/name should fail
func TestCreateDuplicateProjects(t *testing.T) {
  cleanupDB()

  createExampleProject(t)
  res := createExampleProject(t)

  if res.StatusCode != 400 {
    t.Error(fmt.Sprintf("Expected 400 status but got %d", res.StatusCode))
  }
}

// creating a project with invalid json should fail
func TestCreateProjectInvalidJson(t *testing.T) {
	cleanupDB()

	res := createProject(t, `{ "bad" }`)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

	expected := string(body)
	actual := "Invalid JSON\n"
	if expected != actual {
		t.Error(fmt.Sprintf("Expected %s but got %s", expected, actual))
	}
}

// creating a project, then deleting it, should return 200 status and delete the project
func TestDeleteProject(t *testing.T) {
  cleanupDB()

  createExampleProject(t)

  projects := getProjects(t)

  projectId := projects[0].ID
  res := deleteProject(t, projectId)

  if res.StatusCode != 200 {
    t.Errorf("Success expected: Got %d", res.StatusCode)
  }

  newProjects := getProjects(t)

  expected := 0
  actual := len(newProjects)
  if expected != actual {
    t.Error(fmt.Sprintf("Expected %d projects but got %d", expected, actual))
  }
}


// helpers

func createExampleProject(t *testing.T) (*http.Response) {
	return createProject(t, `{"service": "github", "username": "mikec", "name": "example-project"}`)
}

func createProject(t *testing.T, json string) (*http.Response) {
  reader = strings.NewReader(json)

  request, err := http.NewRequest("POST", projectsUrl, reader)
  if err != nil {
    t.Error(err)
  }

  res, err := http.DefaultClient.Do(request)
  if err != nil {
    t.Error(err)
  }
  return res
}

func deleteProject(t *testing.T, id int64) (*http.Response) {
  url := fmt.Sprintf("%s/%d", projectsUrl, id)
  request, err := http.NewRequest("DELETE", url, nil)
  if err != nil {
    t.Error(err)
  }

  res, err := http.DefaultClient.Do(request)
  if err != nil {
    t.Error(err)
  }
  return res
}

func getProjects(t *testing.T) []models.Project {
	request, err := http.NewRequest("GET", projectsUrl, nil)
  if err != nil {
    t.Error(err)
  }

  res, err := http.DefaultClient.Do(request)
  if err != nil {
    t.Error(err)
  }

  decoder := json.NewDecoder(res.Body)
  var projects []models.Project
  err = decoder.Decode(&projects)
  if err != nil {
    t.Error(err)
  }

  return projects
}
