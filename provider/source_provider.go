package provider

type ProjectData struct {
	Url    string `json:"url"`
	Handle string `json:"handle"`
}

// SourceProvider is a service that owns projects with source code,
// such as GitHub, Bitbucket, or Dropbox
type SourceProvider interface {
	AuthProvider

	// GetProjectData returns data for a project owned by a source provider
	GetProjectData(string, string) (*ProjectData, error)

	// GetProjects returns all projects on a source provider owned by a given user
	GetProjects(string, string) ([]ProjectData, error)
}
