package provider

import "fmt"

type ProviderError struct {
	Code          string
	Provider      string
	ProjectHandle string
}

func (e ProviderError) Error() string {
	switch e.Code {
	case "token_invalid":
		return fmt.Sprintf("%s token is invalid", e.Provider)
	case "not_found":
		return fmt.Sprintf("%s project `%s` not found", e.Provider, e.ProjectHandle)
	case "invalid_project_handle":
		return fmt.Sprintf("Invalid %s project handle `%s`", e.Provider, e.ProjectHandle)
	case "get_projects_failed":
		return fmt.Sprintf("Get %s projects failed", e.Provider)
	case "invalid_webook_signature":
		return fmt.Sprintf("Invalid signature in webhook request from %s", e.Provider)
	default:
		return fmt.Sprintf("Unknown %s error", e.Provider)
	}
}

func NewProviderErr(code string, provider string, projectHandle string) *ProviderError {
	return &ProviderError{code, provider, projectHandle}
}

func NewTokenInvalidErr(provider string) *ProviderError {
	return NewProviderErr("token_invalid", provider, "")
}

func NewProjectNotFoundErr(provider string, projectHandle string) *ProviderError {
	return NewProviderErr("not_found", provider, projectHandle)
}

func NewInvalidProjectHandleErr(provider string, projectHandle string) *ProviderError {
	return NewProviderErr("invalid_project_handle", provider, projectHandle)
}

func NewGetProjectsFailedErr(provider string) error {
	return NewProviderErr("get_projects_failed", provider, "")
}

func NewInvalidWebhookSignatureErr(provider string) error {
	return NewProviderErr("invalid_webook_signature", provider, "")
}
