package api_test

import (
	"testing"
)

var runner Runner

func init() {
  runner = Runner{
    apiClient.Projects,
    `{"service": "github", "username": "mikec", "name": "example-1"}`,
    `{"service": "github", "username": "mikec", "name": "example-2"}`,
  }
}

func TestCreateProject(t *testing.T) { runner.RunCreateTest(t) }
func TestGetAllProjects(t *testing.T) { runner.RunGetAllTest(t) }
func TestGetProject(t *testing.T) { runner.RunGetTest(t) }
func TestGetMissingProject(t *testing.T) { runner.RunMissingTest(t) }
func TestCreateDuplicateProject(t *testing.T) { runner.RunDuplicateTest(t) }
func TestCreateProjectWithInvalidJson(t *testing.T) { runner.RunCreateWithInvalidJsonTest(t) }
func TestDeleteProject(t *testing.T) { runner.RunDeleteTest(t) }
