package api_test

import (
  "github.com/mikec/marsupi-api/testutil"

  "github.com/stretchr/testify/assert"

	"testing"
  "io/ioutil"
)

// create project should return 200 status
func TestCreateProject(t *testing.T) {
	cleanupDB()

	res, err := autil.CreateProject(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)
  if err != nil {
    t.Error(err)
  }

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)
}

// create 2 projects, and expect to get 2 projects
func TestGetProjects(t *testing.T) {
	cleanupDB()

	_, err := autil.CreateProject(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)
  if err != nil {
    t.Error(err)
  }

  _, err = autil.CreateProject(`{"service": "github", "username": "mikec", "name": "example-project-2"}`)
  if err != nil {
    t.Error(err)
  }

	projects, _, err := autil.GetProjects()
  if err != nil {
    t.Error(err)
  }

  assert.Len(t, projects, 2, "Wrong number of projects")
}

// create a project and get it by ID
func TestGetProject(t *testing.T) {
  cleanupDB()

  _, err := autil.CreateProject(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)
  if err != nil {
    t.Error(err)
  }

  projects, _, err := autil.GetProjects()
  if err != nil {
    t.Error(err)
  }

  project, _, err := autil.GetProject(projects[0].ID)
  if err != nil {
    t.Error(err)
  }

  assert.Equal(t, project, &projects[0])
}

// creating two projecs with the same service/username/name should fail
func TestCreateDuplicateProjects(t *testing.T) {
  cleanupDB()

  _, err := autil.CreateProject(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)
  if err != nil {
    t.Error(err)
  }

  res, err := autil.CreateProject(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)
  if err != nil {
    t.Error(err)
  }

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
}

// creating a project with invalid json should fail
func TestCreateProjectInvalidJson(t *testing.T) {
	cleanupDB()

	res, err := autil.CreateProject(`{ "bad" }`)
  if err != nil {
    t.Error(err)
  }

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}

  assert.Equal(t, string(body), "Invalid JSON\n")
}

// creating a project, then deleting it, should return 200 status and delete the project
func TestDeleteProject(t *testing.T) {
  cleanupDB()

  _, err := autil.CreateProject(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)
  if err != nil {
    t.Error(err)
  }

  projects, _, err := autil.GetProjects()
  if err != nil {
    t.Error(err)
  }

  res, err := autil.DeleteProject(projects[0].ID)
  if err != nil {
    t.Error(err)
  }

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)

  newProjects, _, err := autil.GetProjects()
  if err != nil {
    t.Error(err)
  }

  assert.Len(t, newProjects, 0)
}

