package api_test

import (
	"testing"

	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/provider"
)

func TestPostProject(t *testing.T) {

	postAuthTest := &PostAuthTest{"mock_seeded_access_token_123"}
	validJson := `{ "project_handle":"jabroni.com/mockers/mock-project", "provider":"jabroni.com" }`

	tests := []*EndpointTest{
		postAuthTest.InvalidJsonEndpointTest(),

		MissingAuthorizationHeaderEndpointTest(validJson),

		InvalidAuthorizationTokenErrorTest(validJson),

		postAuthTest.MissingParamEndpointTest(`{ "project_handle":"asdf" }`, "provider"),

		postAuthTest.MissingParamEndpointTest(`{ "provider":"jabroni.com" }`, "project_handle"),

		postAuthTest.UnsupportedProviderTest(`{ "project_handle":"asdf", "provider":"jnk" }`, "jnk"),

		{
			validJson,
			"a request to save a project for a provider that the user has not authorized",
			"a provider token not found error",
			400,
			api.NewProviderTokenNotFoundErr("jabroni.com"),
			"mock_token_for_user_with_no_provider_tokens",
		},

		{
			`{ "project_handle":"bad_project_handle", "provider":"jabroni.com" }`,
			"a request to save a project that does not exist on the provider",
			"a project not found error",
			400,
			api.NewApiErr(provider.NewProjectDataNotFoundErr("mock/project-x", "jabroni.com").Error()),
			"mock_seeded_access_token_123",
		},

		{
			validJson,
			"a request to save a project that does not exist on the provider",
			"a project not found error",
			400,
			api.NewApiErr(provider.NewProjectDataTokenInvalidErr("mock/project-x", "jabroni.com").Error()),
			"mock_token_for_user_with_bad_jabroni.com_token",
		},

		{
			validJson,
			"a request to save a valid project",
			"success",
			200,
			api.ApiResponse{"success!"},
			"mock_seeded_access_token_123",
		},
	}

	RunEndpointTests(t, "POST", "projects", tests)

}
