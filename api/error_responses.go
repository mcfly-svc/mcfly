package api

import (
	"fmt"
)

type ApiError struct {
	Error string `json:"error"`
}

func NewApiErr(errorMessage string) *ApiError {
	return &ApiError{errorMessage}
}

func NewServerErr() *ApiError {
	return NewApiErr("Unknown server error. That's bad!")
}

func NewInvalidJsonErr() *ApiError {
	return NewApiErr("Invalid JSON")
}

func NewInvalidParamErr(name string, val string) *ApiError {
	return NewApiErr(fmt.Sprintf("Invalid parameter %s=%s", name, val))
}

func NewMissingParamErr(param string) *ApiError {
	return NewApiErr(fmt.Sprintf("Missing %s parameter", param))
}

func NewInvalidRequestParamsErr() *ApiError {
	return NewApiErr("Invalid request parameters")
}

func NewDuplicateCreateErr(entityName string) *ApiError {
	return NewApiErr(fmt.Sprintf("Could not create %s because it already exists", entityName))
}

func NewCreateFailedErr(entityName string) *ApiError {
	return NewApiErr(fmt.Sprintf("Failed to create %s", entityName))
}

func NewGetEntitiesErr(entitiesName string) *ApiError {
	return NewApiErr(fmt.Sprintf("Failed to get %s", entitiesName))
}

func NewGetEntityErr(entityName string, ID int64) *ApiError {
	return NewApiErr(fmt.Sprintf("Failed to get %s where id=%d", entityName, ID))
}

func NewDeleteFailedErr(entityName string) *ApiError {
	return NewApiErr(fmt.Sprintf("Failed to delete %s", entityName))
}

func NewUnsupportedProviderErr(provider string) *ApiError {
	return NewApiErr(fmt.Sprintf("Unsupported provider %s", provider))
}

func NewInvalidTokenErr(token string) *ApiError {
	return NewApiErr(fmt.Sprintf("Invalid %s token", token))
}

func NewAuthorizationHeaderRequiredErr() *ApiError {
	return NewApiErr("Authorization header required")
}

func NewInvalidAuthTokenError(token string) *ApiError {
	return NewApiErr(fmt.Sprintf("Auth token %s is not valid", token))
}
