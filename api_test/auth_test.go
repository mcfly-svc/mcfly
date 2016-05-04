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
	}

	RunEndpointTests(t, "POST", "login", tests)
}
