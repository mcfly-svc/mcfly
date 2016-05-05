package api_test

import (
	"testing"

	"github.com/mikec/msplapi/api"
)

func TestLogin(t *testing.T) {

	tests := []*EndpointTest{
		InvalidJsonEndpointTest(),

		MissingParamEndpointTest(`{ "token":"abc123" }`, "token_type"),

		MissingParamEndpointTest(`{ "token_type":"jabroni.com" }`, "token"),

		{
			`{ "token":"abc123", "token_type":"junk-service" }`,
			"an unsupported token type",
			"an unsupported token type error",
			400,
			api.NewUnsupportedTokenTypeErr("junk-service"),
		},

		{
			`{ "token":"badtoken", "token_type":"jabroni.com" }`,
			"an invalid token",
			"an invalid token error",
			400,
			api.NewInvalidTokenErr("jabroni.com"),
		},

		{
			`{ "token":"mock_jabroni.com_token_123", "token_type":"jabroni.com"}`,
			"a valid provider token owned by an existing user",
			"expected user data with the existing access token",
			200,
			api.LoginResp{Name: "Matt Mockman", AccessToken: "mock_seeded_access_token_123"},
			/*map[string]interface{}{
				"name":         "Matt Mockman",
				"access_token": "mock_seeded_access_token_123",
			},*/
		},

		{
			`{ "token":"mock_jb_token_123", "token_type":"jabroni.com"}`,
			"a valid provider token owned by a user that does not exist",
			"expected user data with a newly generated access token",
			200,
			api.LoginResp{Name: "Mike Jimmers", AccessToken: "mock_generated_access_token_123"},
			/*map[string]interface{}{
				"name":         "Mike Jimmers",
				"access_token": "mock_generated_access_token_123",
			},*/
		},
	}

	RunEndpointTests(t, "POST", "login", tests)
}
