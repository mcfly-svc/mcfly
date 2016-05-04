package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/mikec/msplapi/api"
	. "github.com/smartystreets/goconvey/convey"
)

type ApiErrorTest struct {
	*testing.T
	EndpointPath   string
	JSON           string
	Desc           string
	Should         string
	ExpectApiError *api.ApiError
}

func RunPostInvalidJsonTest(t *testing.T, endpointPath string) {
	RunPostErrorTest(&ApiErrorTest{
		t,
		endpointPath,
		`{"jnk"}`,
		"invalid JSON",
		"Invalid JSON error",
		api.NewInvalidJsonErr(),
	})
}

func RunMissingPostParamTest(t *testing.T, endpointPath string, json string, paramName string) {
	RunPostErrorTest(&ApiErrorTest{
		t,
		endpointPath,
		json,
		fmt.Sprintf("missing `%s` parameter", paramName),
		fmt.Sprintf("missing `%s` parameter", paramName),
		api.NewMissingParamErr(paramName),
	})
}

func RunPostErrorTest(at *ApiErrorTest) {

	Convey(fmt.Sprintf("POST to /%s endpoint with %s", at.EndpointPath, at.Desc), at, func() {

		res, err := apiClient.Context.DoPost(at.EndpointPath, &at.JSON, nil)
		if err != nil {
			at.Error(err)
		}

		Convey("Should respond with Status 400", func() {
			So(res.StatusCode, ShouldEqual, 400)
		})

		Convey(fmt.Sprintf("Should respond with %s", at.Should), func() {
			SoErrorRes(at.T, res, at.ExpectApiError)
		})

	})

}

func SoErrorRes(t *testing.T, res *http.Response, expectedApiError *api.ApiError) {
	apiErr, err := errorFromRes(res)
	if err != nil {
		t.Error(err)
	}
	So(apiErr, ShouldResemble, expectedApiError)
}

func dataFromRes(res *http.Response, v interface{}) ([]byte, error) {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func errorFromRes(res *http.Response) (*api.ApiError, error) {
	apiErr := api.ApiError{}
	_, err := dataFromRes(res, &apiErr)
	if err != nil {
		return nil, err
	}
	return &apiErr, nil
}
