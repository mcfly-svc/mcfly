package api_test

import (
	"testing"
)

var tests EndpointTests

func init() {
  tests = EndpointTests{
    autil.Projects,
    `{"service": "github", "username": "mikec", "name": "example-1"}`,
    `{"service": "github", "username": "mikec", "name": "example-2"}`,
  }
}

func TestCreateProject(t *testing.T) { tests.RunCreateTest(t) }
func TestGetProjects(t *testing.T) { tests.RunGetAllTest(t) }
func TestGetProject(t *testing.T) { tests.RunGetTest(t) }
func TestMissingProject(t *testing.T) { tests.RunMissingTest(t) }
func TestDuplicateProjects(t *testing.T) { tests.RunDuplicateTest(t) }
func TestCreateProjectWithInvalidJson(t *testing.T) { tests.RunCreateWithInvalidJsonTest(t) }
func TestDeleteProject(t *testing.T) { tests.RunDeleteTest(t) }
