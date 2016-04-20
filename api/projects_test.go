package api_test

import (
  "github.com/mikec/marsupi-api/models"
  "github.com/mikec/marsupi-api/api"
  "github.com/mikec/marsupi-api/testutil"

  "github.com/stretchr/testify/assert"

	"testing"
)

// create project should return 200 status
func TestCreateProject(t *testing.T) {
	cleanupDB()

  et := EndpointTester{t, autil.Projects}
  res := et.Create(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)
}

// create 2 projects, and expect to get 2 projects
func TestGetProjects(t *testing.T) {
	cleanupDB()

  et := EndpointTester{t, autil.Projects}
  et.Create(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)
  et.Create(`{"service": "github", "username": "mikec", "name": "example-project-2"}`)
  projects := et.GetAll()

  assert.Len(t, projects, 2, "Wrong number of projects")
}


// create a project and get it by ID
func TestGetProject(t *testing.T) {
  cleanupDB()

  et := EndpointTester{t, autil.Projects}
  et.Create(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)
  projects := et.GetAll()
  p1 := projects[0].(models.Project)
  p2 := et.Get(p1.ID)

  assert.Equal(t, p2, p1)
}


// get a project that doesn't exist
func TestGetProjectThatDoesNotExist(t *testing.T) {
  cleanupDB()

  et := EndpointTester{t, autil.Projects}
  res := et.GetRes(123)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
  rt.ExpectResponseBody(api.ApiError{"Failed to get project with id=123"})

}

// creating two projecs with the same service/username/name should fail
func TestCreateDuplicateProjects(t *testing.T) {
  cleanupDB()

  et := EndpointTester{t, autil.Projects}
  et.Create(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)
  res := et.Create(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
}

// creating a project with invalid json should fail
func TestCreateProjectInvalidJson(t *testing.T) {
	cleanupDB()

  et := EndpointTester{t, autil.Projects}
	res := et.Create(`{ "bad" }`)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(400)
  rt.ExpectResponseBody(api.InvalidJsonApiErr)
}


// creating a project, then deleting it, should return 200 status and delete the project
func TestDeleteProject(t *testing.T) {
  cleanupDB()

  et := EndpointTester{t, autil.Projects}
  et.Create(`{"service": "github", "username": "mikec", "name": "example-project-1"}`)
  projects := et.GetAll()
  p := projects[0].(models.Project)
  res := et.Delete(p.ID)

  rt := &testutil.ResponseTest{t, res}
  rt.ExpectHttpStatus(200)

  projects = et.GetAll()

  assert.Len(t, projects, 0)
}

