package api_test

import (
	"testing"

	"github.com/mikec/msplapi/api"
)

func TestPostProject(t *testing.T) {

	validJson := `{ "project_name":"jabroni.com/mockers/mock-project", "provider":"jabroni.com" }`

	tests := []*EndpointTest{
		InvalidJsonEndpointTest(),

		MissingAuthorizationHeaderEndpointTest(validJson),

		InvalidAuthorizationTokenErrorTest(validJson),

		MissingParamEndpointTest(`{ "project_name":"asdf" }`, "provider"),

		MissingParamEndpointTest(`{ "provider":"jabroni.com" }`, "project_name"),

		{
			validJson,
			"a request to save a project for a provider that the user has not authorized",
			"a provider token not found error",
			400,
			api.NewProviderTokenNotFoundErr("jabroni.com"),
			"mock_token_for_user_with_no_provider_tokens",
		},
	}

	RunEndpointTests(t, "POST", "projects", tests)

}
