package api_test

import (
	"testing"

	"github.com/mcfly-svc/mcfly/api"
	"github.com/mcfly-svc/mcfly/models"
)

func TestSaveBuild(t *testing.T) {
	validJson := `{"handle":"abc", "project_handle":"mattmocks/project-2", 
				"provider":"jabroni.com", "provider_url":"https://jabroni.com/mock-builds/abc"}`

	cleanupDB()
	RunEndpointTests(t, "POST", "builds", []*EndpointTest{
		{
			validJson,
			"a request to save a build",
			"success",
			200,
			api.NewSuccessResponse(),
			"",
		},
	})
	p := models.Project{Handle: "mattmocks/project-2", SourceProvider: "jabroni.com"}
	commits := make(map[string]bool)
	commits["abc"] = true
	assertProjectHasBuildCommits(t, "When a build is saved for a project", &p, commits)
}
