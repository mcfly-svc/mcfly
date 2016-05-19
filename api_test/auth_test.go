package api_test

import (
	"testing"

	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/api/apidata"
)

func TestLogin(t *testing.T) {

	afterAuthTest := &AfterAuthTest{"mock_seeded_access_token_123"}

	tests := []*EndpointTest{
		afterAuthTest.InvalidJsonEndpointTest(),

		afterAuthTest.MissingParamEndpointTest(`{ "token":"abc123" }`, "provider"),

		afterAuthTest.MissingParamEndpointTest(`{ "provider":"jabroni.com" }`, "token"),

		afterAuthTest.UnsupportedProviderTest(`{ "token":"abc123", "provider":"junk-service" }`, "junk-service"),

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
			apidata.LoginResp{Name: strPtr("Matt Mockman"), AccessToken: "mock_seeded_access_token_123"},
			"",
		},

		{
			`{ "token":"mock_dne_user_token_123", "provider":"jabroni.com"}`,
			"a valid provider token owned by a user that does not exist",
			"expected user data with a newly generated access token",
			200,
			apidata.LoginResp{Name: strPtr("Mike Jimmers"), AccessToken: "mock_generated_access_token_123"},
			"",
		},
	}

	RunEndpointTests(t, "POST", "login", tests)
}
