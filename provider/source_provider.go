package provider

import (
	"net/http"
	"strings"
)

type ProjectData struct {
	Url    string `json:"url"`
	Handle string `json:"handle"`
}

type ProjectUpdateData struct {
	ProjectHandle string
	Builds        []string
}

// SourceProvider is a service that owns projects with source code,
// such as GitHub, Bitbucket, or Dropbox
type SourceProvider interface {
	AuthProvider

	// GetProjectData returns data for a project owned by a source provider
	GetProjectData(string, string) (*ProjectData, error)

	// GetProjects returns all projects on a source provider owned by a given user
	GetProjects(string, string) ([]ProjectData, error)

	// CreateProjectUpdateHook creates a webhook on the source provider that will notify
	// msplapi when a given project is updated
	CreateProjectUpdateHook(string, string) error

	// DecodeProjectUpdateRequest decodes the request made to the ProjectUpdate webhook
	// by the source provider, and returns an array of builds that were added in the
	// project update
	DecodeProjectUpdateRequest(*http.Request) (*ProjectUpdateData, error)
}

type SourceProviderConfig struct {
	ProjectUpdateHookUrlFmt string
}

func (self *SourceProviderConfig) GetProjectUpdateHookUrl(key string) string {
	return strings.Replace(self.ProjectUpdateHookUrlFmt, "{provider}", key, 1)
}
