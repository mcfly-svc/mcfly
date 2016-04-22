package api_test

import (
	"testing"
)

var projectsRunner Runner

func init() {
  projectsRunner = Runner{
    apiClient.Projects,
    `{"service": "github", "username": "mikec", "name": "example-1"}`,
    `{"service": "github", "username": "mikec", "name": "example-2"}`,
  }
}

func TestCreateProject(t *testing.T) { projectsRunner.RunCreateTest(t) }
func TestGetAllProjects(t *testing.T) { projectsRunner.RunGetAllTest(t) }
func TestGetProject(t *testing.T) { projectsRunner.RunGetTest(t) }
func TestGetMissingProject(t *testing.T) { projectsRunner.RunMissingTest(t) }
func TestCreateDuplicateProject(t *testing.T) { projectsRunner.RunDuplicateTest(t) }
func TestCreateProjectWithInvalidJson(t *testing.T) { projectsRunner.RunCreateWithInvalidJsonTest(t) }
func TestDeleteProject(t *testing.T) { projectsRunner.RunDeleteTest(t) }
func TestGetProjectInvalidID(t *testing.T) { projectsRunner.RunInvalidGetTest(t) }
