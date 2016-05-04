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

type EndpointTest struct {
	JSON         string
	Desc         string
	Should       string
	ExpectStatus int
	ExpectBody   interface{}
}

func InvalidJsonEndpointTest() *EndpointTest {
	return &EndpointTest{
		`{"jnk"}`,
		"invalid JSON",
		"Invalid JSON error",
		400,
		api.NewInvalidJsonErr(),
	}
}

func MissingParamEndpointTest(json string, paramName string) *EndpointTest {
	return &EndpointTest{
		json,
		fmt.Sprintf("missing `%s` parameter", paramName),
		fmt.Sprintf("missing `%s` parameter error", paramName),
		400,
		api.NewMissingParamErr(paramName),
	}
}

func RunEndpointTests(t *testing.T, httpMethod string, endpointPath string, tests []*EndpointTest) {
	for _, et := range tests {
		Convey(fmt.Sprintf("%s to /%s endpoint %s", httpMethod, endpointPath, et.Desc), t, func() {

			var jsonData *string
			if et.JSON == "" {
				jsonData = nil
			} else {
				jsonData = &et.JSON
			}

			res, err := apiClient.Context.DoReq(httpMethod, endpointPath, jsonData, nil)
			if err != nil {
				t.Error(err)
			}

			Convey(fmt.Sprintf("Should respond with Status %d", et.ExpectStatus), func() {
				So(res.StatusCode, ShouldEqual, et.ExpectStatus)
			})

			Convey(fmt.Sprintf("Should respond with %s", et.Should), func() {
				So(resBody(t, res), ShouldEqual, marshalJson(t, et.ExpectBody))
			})

		})
	}
}

func resBody(t *testing.T, res *http.Response) string {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	return string(b)
}

func marshalJson(t *testing.T, v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		t.Error(err)
	}
	return string(b)
}
