package api_test

import (
	"testing"

	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/provider"
)

func TestPostProject(t *testing.T) {

	afterAuthTest := &AfterAuthTest{"mock_seeded_access_token_123"}
	validJson := `{ "project_handle":"mock/project-x", "provider":"jabroni.com" }`

	tests := []*EndpointTest{
		afterAuthTest.InvalidJsonEndpointTest(),

		MissingAuthorizationHeaderEndpointTest(validJson),

		InvalidAuthorizationTokenErrorTest(validJson),

		afterAuthTest.MissingParamEndpointTest(`{ "project_handle":"asdf" }`, "provider"),

		afterAuthTest.MissingParamEndpointTest(`{ "provider":"jabroni.com" }`, "project_handle"),

		afterAuthTest.UnsupportedProviderTest(`{ "project_handle":"asdf", "provider":"jnk" }`, "jnk"),

		{
			validJson,
			"a request to save a project for a provider that the user has not authorized",
			"a provider token not found error",
			400,
			api.NewProviderTokenNotFoundErr("jabroni.com"),
			"mock_token_for_user_with_no_provider_tokens",
		},

		{
			`{ "project_handle":"project_handle_dne", "provider":"jabroni.com" }`,
			"a request to save a project that does not exist on the provider",
			"a project not found error",
			400,
			api.NewApiErr(provider.NewProjectDataNotFoundErr("mock/project-x", "jabroni.com").Error()),
			"mock_seeded_access_token_123",
		},

		{
			`{ "project_handle":"invalid_project_handle", "provider":"jabroni.com" }`,
			"a request to save a project with an invalid project handle",
			"a project handle invalid error",
			400,
			api.NewApiErr(provider.NewProjectDataInvalidHandleErr("invalid_project_handle", "jabroni.com").Error()),
			"mock_seeded_access_token_123",
		},

		{
			validJson,
			"a request to save a project when the saved provider token is invalid",
			"a token invalid error",
			400,
			api.NewApiErr(provider.NewProviderTokenInvalidErr("jabroni.com").Error()),
			"mock_token_for_user_with_bad_jabroni.com_token",
		},

		{
			validJson,
			"a request to save a valid project",
			"success",
			200,
			api.PostProjectResp{"mock/project-x", "https://jabroni.com/mock/project-x", "jabroni.com"},
			"mock_seeded_access_token_123",
		},
	}

	RunEndpointTests(t, "POST", "projects", tests)

}

func TestGetProviderProjects(t *testing.T) {
	afterAuthTest := &AfterAuthTest{"mock_seeded_access_token_123"}

	tests := []*EndpointTest{
		MissingAuthorizationHeaderEndpointTest(""),

		InvalidAuthorizationTokenErrorTest(""),
	}
	RunEndpointTests(t, "GET", "jabroni.com/projects", tests)

	invalidProviderTests := []*EndpointTest{
		afterAuthTest.UnsupportedProviderTest("", "jnk"),
	}
	RunEndpointTests(t, "GET", "jnk/projects", invalidProviderTests)
}
