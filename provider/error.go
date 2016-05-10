package provider

import "fmt"

type ProviderTokenInvalidError struct {
	Provider string
}

func (e ProviderTokenInvalidError) Error() string {
	return fmt.Sprintf("%s token is invalid", e.Provider)
}

func NewProviderTokenInvalidErr(token string) *ProviderTokenInvalidError {
	return &ProviderTokenInvalidError{token}
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
	} else if e.Name == "invalid_handle" {
		return fmt.Sprintf("%s. Invalid %s project handle `%s`", msg, e.Provider, e.ProjectHandle)
	}
	return msg
}

func NewProjectDataNotFoundErr(projectHandle string, provider string) *GetProjectDataError {
	return &GetProjectDataError{"not_found", projectHandle, provider}
}

func NewProjectDataInvalidHandleErr(projectHandle string, provider string) *GetProjectDataError {
	return &GetProjectDataError{"invalid_handle", projectHandle, provider}
}
