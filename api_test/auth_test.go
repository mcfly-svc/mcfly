package api_test

import (
	"testing"

	"github.com/mikec/msplapi/api"
)

func TestLogin(t *testing.T) {

	postAuthTest := &PostAuthTest{"mock_seeded_access_token_123"}

	tests := []*EndpointTest{
		postAuthTest.InvalidJsonEndpointTest(),

		postAuthTest.MissingParamEndpointTest(`{ "token":"abc123" }`, "provider"),

		postAuthTest.MissingParamEndpointTest(`{ "provider":"jabroni.com" }`, "token"),

		postAuthTest.UnsupportedProviderTest(`{ "token":"abc123", "provider":"junk-service" }`, "junk-service"),

		{
			`{ "token":"badtoken", "provider":"jabroni.com" }`,
			"an invalid token",
			"an invalid token error",
			400,
			api.NewInvalidTokenErr("jabroni.com"),
			"",
		},

		{
			`{ "token":"mock_jabroni.com_token_123", "provider":"jabroni.com"}`,
			"a valid provider token owned by an existing user",
			"expected user data with the existing access token",
			200,
			api.LoginResp{Name: "Matt Mockman", AccessToken: "mock_seeded_access_token_123"},
			/*map[string]interface{}{
				"name":         "Matt Mockman",
				"access_token": "mock_seeded_access_token_123",
			},*/
			"",
		},

		{
			`{ "token":"mock_jb_token_123", "provider":"jabroni.com"}`,
			"a valid provider token owned by a user that does not exist",
			"expected user data with a newly generated access token",
			200,
			api.LoginResp{Name: "Mike Jimmers", AccessToken: "mock_generated_access_token_123"},
			/*map[string]interface{}{
				"name":         "Mike Jimmers",
				"access_token": "mock_generated_access_token_123",
			},*/
			"",
		},
	}

	RunEndpointTests(t, "POST", "login", tests)
}
