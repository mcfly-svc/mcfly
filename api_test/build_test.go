package api_test

import (
	"testing"

	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/models"
)

func TestSaveBuild(t *testing.T) {
	cleanupDB()
	RunEndpointTests(t, "POST", "builds", []*EndpointTest{
		{
			`{"handle":"abc", "project_handle":"mattmocks/project-3", "provider":"jabroni.com"}`,
			"a request to save a build",
			"success",
			200,
			api.NewSuccessResponse(),
			"",
		},
	})
	p := models.Project{Handle: "mattmocks/project-3", SourceProvider: "jabroni.com"}
	commits := make(map[string]bool)
	commits["abc"] = true
	assertProjectHasBuildCommits(t, "When a build is saved for a project", &p, commits)
}
