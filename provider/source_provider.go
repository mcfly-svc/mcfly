package provider

import "fmt"

type ProjectData struct {
	Url string
}

// SourceProvider is a service that owns projects with source code,
// such as GitHub, Bitbucket, or Dropbox
type SourceProvider interface {
	AuthProvider

	// GetProjectData returns data for a project owned by a source provider
	GetProjectData(string, string) (*ProjectData, error)
}

type GetProjectDataError struct {
	Name          string
	ProjectHandle string
	Provider      string
}

func (e GetProjectDataError) Error() string {
	msg := fmt.Sprintf("Unable to get project from %s", e.Provider)
	if e.Name == "not_found" {
		return fmt.Sprintf("%s. %s project `%s` not found", msg, e.Provider, e.ProjectHandle)
	} else if e.Name == "invalid_token" {
		return fmt.Sprintf("%s. %s token is invalid", msg, e.Provider)
	} else if e.Name == "invalid_handle" {
		return fmt.Sprintf("%s. Invalid %s project handle `%s`", msg, e.Provider, e.ProjectHandle)
	}
	return msg
}

func NewProjectDataNotFoundErr(projectHandle string, provider string) *GetProjectDataError {
	return &GetProjectDataError{"not_found", projectHandle, provider}
}

func NewProjectDataTokenInvalidErr(projectHandle string, provider string) *GetProjectDataError {
	return &GetProjectDataError{"invalid_token", projectHandle, provider}
}

func NewProjectDataInvalidHandleErr(projectHandle string, provider string) *GetProjectDataError {
	return &GetProjectDataError{"invalid_handle", projectHandle, provider}
}
