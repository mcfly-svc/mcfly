package api

import (
	"fmt"
	"reflect"

	"gopkg.in/validator.v2"

	"github.com/mikec/msplapi/models"
)

// ValidateAuthorization checks the "Authorization" header in the request. If this header doesn't
// exist or contains an invalid token, it responds with an ApiError and returns nil. If the
// auth token is valid, it responds with the User data corresponding with that auth token
func (r *Responder) ValidateAuthorization(db models.Datastore) *models.User {
	authHeader := r.Request.Header["Authorization"]
	if len(authHeader) == 0 {
		r.RespondWithError(NewAuthorizationHeaderRequiredErr())
		return nil
	}

	authToken := authHeader[0]

	user, err := db.GetUserByAccessToken(authToken)
	if err != nil {
		r.RespondWithServerError(err)
		return nil
	}

	if user == nil {
		r.RespondWithError(NewInvalidAuthTokenError(authToken))
		return nil
	}

	return user
}

// ValidateRequestData uses gopkg.in/validator.v2 to validate request data. It returns the appropriate
// ApiError based on the result of the validation, or nil if there are no validation errors
func (r *Responder) ValidateRequestData(reqData interface{}) bool {
	if reflect.TypeOf(reqData).Kind() != reflect.Ptr {
		r.RespondWithServerError(fmt.Errorf("reqData argument to ValidateRequestData must be a pointer"))
		return false
	}
	if err := validator.Validate(reqData); err != nil {
		errs := err.(validator.ErrorMap)
		var badParam *string
		val := reflect.ValueOf(reqData).Elem()
		for i := 0; i < val.NumField(); i++ {
			f := val.Type().Field(i)
			if len(errs[f.Name]) > 0 {
				tagName := f.Tag.Get("json")
				badParam = &tagName
				break
			}
		}
		if badParam != nil {
			r.RespondWithError(NewMissingParamErr(*badParam))
		} else {
			r.RespondWithError(NewInvalidRequestParamsErr())
		}
		return false
	}
	return true
}
