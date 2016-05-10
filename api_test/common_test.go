package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/client"
	. "github.com/smartystreets/goconvey/convey"
)

type EndpointTest struct {
	JSON         string
	Desc         string
	Should       string
	ExpectStatus int
	ExpectBody   interface{}
	AccessToken  string
}

func MissingAuthorizationHeaderEndpointTest(json string) *EndpointTest {
	return &EndpointTest{
		json,
		"nothing in the Authorization header",
		"an authorization header required error",
		400,
		api.NewAuthorizationHeaderRequiredErr(),
		"",
	}
}

func InvalidAuthorizationTokenErrorTest(json string) *EndpointTest {
	return &EndpointTest{
		json,
		"an invalid access token",
		"an invalid access token error",
		400,
		api.NewInvalidAuthTokenError("mock_invalid_access_token"),
		"mock_invalid_access_token",
	}
}

type AfterAuthTest struct {
	AccessToken string
}

func (t *AfterAuthTest) MissingParamEndpointTest(json string, paramName string) *EndpointTest {
	return &EndpointTest{
		json,
		fmt.Sprintf("missing `%s` parameter", paramName),
		fmt.Sprintf("missing `%s` parameter error", paramName),
		400,
		api.NewMissingParamErr(paramName),
		t.AccessToken,
	}
}

func (t *AfterAuthTest) InvalidJsonEndpointTest() *EndpointTest {
	return &EndpointTest{
		`{"jnk"}`,
		"invalid JSON",
		"Invalid JSON error",
		400,
		api.NewInvalidJsonErr(),
		t.AccessToken,
	}
}

func (t *AfterAuthTest) UnsupportedProviderTest(json string, badProvider string) *EndpointTest {
	return &EndpointTest{
		json,
		"an unsupported provider",
		"an unsupported provider error",
		400,
		api.NewUnsupportedProviderErr(badProvider),
		t.AccessToken,
	}
}

func RunEndpointTests(t *testing.T, httpMethod string, endpointPath string, tests []*EndpointTest) {
	for _, et := range tests {
		Convey(fmt.Sprintf("%s to /%s endpoint %s", httpMethod, endpointPath, et.Desc), t, func() {
			resetDB()

			var jsonData *string
			if et.JSON == "" {
				jsonData = nil
			} else {
				jsonData = &et.JSON
			}

			reqOpts := &client.RequestOptions{}
			if et.AccessToken != "" {
				reqOpts.AuthHeader = &et.AccessToken
			}

			res, err := apiClient.DoReq(httpMethod, endpointPath, jsonData, reqOpts)
			if err != nil {
				t.Error(err)
			}

			Convey(fmt.Sprintf("Should respond with Status %d", et.ExpectStatus), func() {
				So(res.StatusCode, ShouldEqual, et.ExpectStatus)
			})

			Convey(fmt.Sprintf("Should respond with %s", et.Should), func() {
				switch v := et.ExpectBody.(type) {
				case map[string]interface{}:
					soFieldsShouldEqual(t, res, v)
				default:
					soBodyShouldEqual(t, res, v)
				}
			})

		})
	}
}

func soFieldsShouldEqual(t *testing.T, res *http.Response, expectedFields map[string]interface{}) {
	actualFields := resBodyMap(t, res)
	for f, v := range expectedFields {
		So(actualFields[f], ShouldEqual, v)
	}
}

func soBodyShouldEqual(t *testing.T, res *http.Response, v interface{}) {
	So(resBody(t, res), ShouldEqual, marshalJson(t, v))
}

func resBody(t *testing.T, res *http.Response) string {
	return string(resBodyBytes(t, res))
}

func resBodyBytes(t *testing.T, res *http.Response) []byte {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	return b
}

func resBodyMap(t *testing.T, res *http.Response) map[string]interface{} {
	b := resBodyBytes(t, res)
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		t.Error(err)
	}
	return m
}

func marshalJson(t *testing.T, v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		t.Error(err)
	}
	return string(b)
}
