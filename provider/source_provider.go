package provider

import (
	"fmt"
	"net/http"
	"strings"
)

type ProjectData struct {
	Url    string `json:"url"`
	Handle string `json:"handle"`
}

type BuildData struct {
	Url    *string
	Handle string
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
	// project update. It should return an error if the request did not originate from
	// the provider
	DecodeProjectUpdateRequest(*http.Request) (*ProjectUpdateData, error)

	// GetBuildData gets data for a given build
	GetBuildData(string, string, string) (*BuildData, error)

	// GetBuildConfig gets the config (mspl.json) for a given build
	GetBuildConfig(string, string, string) (*BuildConfig, error)
}

type SourceProviderConfig struct {
	ProjectUpdateHookUrlFmt string
	WebhookSecret           string
}

func (self *SourceProviderConfig) GetProjectUpdateHookUrl(key string) string {
	return strings.Replace(self.ProjectUpdateHookUrlFmt, "{provider}", key, 1)
}

func GetSourceProviders(config *SourceProviderConfig) map[string]SourceProvider {
	sourceProviders := make(map[string]SourceProvider)
	github := GitHub{GitHubClient: &GoGitHubClient{}, SourceProviderConfig: config}
	dropbox := Dropbox{SourceProviderConfig: config}
	sourceProviders[github.Key()] = &github
	sourceProviders[dropbox.Key()] = &dropbox
	return sourceProviders
}

func GetSourceProvider(key string, config *SourceProviderConfig) (SourceProvider, error) {
	sourceProviders := GetSourceProviders(config)
	sp := sourceProviders[key]
	if sp == nil {
		return nil, NewSourceProviderNotSupportedErr(key)
	}
	return sp, nil
}

func NewSourceProviderNotSupportedErr(key string) error {
	return fmt.Errorf("Source provider %s not supported", key)
}
