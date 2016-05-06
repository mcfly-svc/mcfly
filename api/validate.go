package api

import (
	"fmt"
	"net/http"
	"reflect"

	"gopkg.in/validator.v2"

	"github.com/mikec/msplapi/models"
)

func ValidateAuthorization(db models.Datastore, req *http.Request) (*models.User, *ApiError, error) {
	authHeader := req.Header["Authorization"]
	if len(authHeader) == 0 {
		return nil, NewAuthorizationHeaderRequiredErr(), nil
	}

	authToken := authHeader[0]

	user, err := db.GetUserByAccessToken(authToken)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, NewInvalidAuthTokenError(authToken), nil
	}

	return user, nil, nil
}

// ValidateRequestData uses gopkg.in/validator.v2 to validate request data. It returns the appropriate
// ApiError based on the result of the validation, or nil if there are no validation errors
func ValidateRequestData(reqData interface{}) *ApiError {
	if reflect.TypeOf(reqData).Kind() != reflect.Ptr {
		panic(fmt.Errorf("reqData argument to ValidateRequestData must be a pointer"))
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
			return NewMissingParamErr(*badParam)
		} else {
			return NewInvalidRequestParamsErr()
		}
	}
	return nil
}
