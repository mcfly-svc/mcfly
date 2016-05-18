package api

import (
	"fmt"

	"github.com/mikec/msplapi/api/apidata"
)

func NewApiErr(errorMessage string) *apidata.ApiError {
	return &apidata.ApiError{errorMessage}
}

func NewServerErr() *apidata.ApiError {
	return NewApiErr("Unknown server error. That's bad!")
}

func NewInvalidJsonErr() *apidata.ApiError {
	return NewApiErr("Invalid JSON")
}

func NewInvalidParamErr(name string, val string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("Invalid parameter %s=%s", name, val))
}

func NewMissingParamErr(param string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("Missing %s parameter", param))
}

func NewInvalidRequestParamsErr() *apidata.ApiError {
	return NewApiErr("Invalid request parameters")
}

func NewDuplicateCreateErr(entityName string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("Could not create %s because it already exists", entityName))
}

func NewCreateFailedErr(entityName string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("Failed to create %s", entityName))
}

func NewGetEntitiesErr(entitiesName string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("Failed to get %s", entitiesName))
}

func NewGetEntityErr(entityName string, ID int64) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("Failed to get %s where id=%d", entityName, ID))
}

func NewDeleteFailedErr(entityName string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("Failed to delete %s", entityName))
}

func NewUnsupportedProviderErr(provider string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("Unsupported provider %s", provider))
}

func NewInvalidTokenErr(token string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("Invalid %s token", token))
}

func NewAuthorizationHeaderRequiredErr() *apidata.ApiError {
	return NewApiErr("Authorization header required")
}

func NewInvalidAuthTokenError(token string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("Auth token %s is not valid", token))
}

func NewProviderTokenNotFoundErr(provider string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("%s token not found", provider))
}

func NewNotFoundErr(entity string, handle string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("%s %s not found", entity, handle))
}

func NewDuplicateErr(entity string, handle string) *apidata.ApiError {
	return NewApiErr(fmt.Sprintf("%s %s already added", entity, handle))
}
