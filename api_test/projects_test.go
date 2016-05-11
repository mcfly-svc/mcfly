package api_test

import (
	"testing"

	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/provider"
)

func TestPostProject(t *testing.T) {

	afterAuthTest := &AfterAuthTest{"mock_seeded_access_token_123"}
	validJson := `{ "handle":"mock/project-x", "provider":"jabroni.com" }`

	tests := []*EndpointTest{
		afterAuthTest.InvalidJsonEndpointTest(),

		MissingAuthorizationHeaderEndpointTest(validJson),

		InvalidAuthorizationTokenErrorTest(validJson),

		afterAuthTest.MissingParamEndpointTest(`{ "handle":"asdf" }`, "provider"),

		afterAuthTest.MissingParamEndpointTest(`{ "provider":"jabroni.com" }`, "handle"),

		afterAuthTest.UnsupportedProviderTest(`{ "handle":"asdf", "provider":"jnk" }`, "jnk"),

		{
			validJson,
			"a request to save a project for a provider that the user has not authorized",
			"a provider token not found error",
			400,
			api.NewProviderTokenNotFoundErr("jabroni.com"),
			"mock_token_for_user_with_no_provider_tokens",
		},

		{
			`{ "handle":"project_handle_dne", "provider":"jabroni.com" }`,
			"a request to save a project that does not exist on the provider",
			"a project not found error",
			400,
			api.NewApiErr(provider.NewProjectNotFoundErr("jabroni.com", "mock/project-x").Error()),
			"mock_seeded_access_token_123",
		},

		{
			`{ "handle":"invalid_project_handle", "provider":"jabroni.com" }`,
			"a request to save a project with an invalid project handle",
			"a project handle invalid error",
			400,
			api.NewApiErr(provider.NewInvalidProjectHandleErr("jabroni.com", "invalid_project_handle").Error()),
			"mock_seeded_access_token_123",
		},

		{
			validJson,
			"a request to save a project when the saved provider token is invalid",
			"a token invalid error",
			400,
			api.NewApiErr(provider.NewTokenInvalidErr("jabroni.com").Error()),
			"mock_token_for_user_with_bad_jabroni.com_token",
		},

		{
			`{ "handle":"mattmocks/project-1", "provider":"jabroni.com" }`,
			"a request to save a project that already exists",
			"a project exists error",
			400,
			api.NewDuplicateErr("project", "mattmocks/project-1"),
			"mock_seeded_access_token_123",
		},

		{
			validJson,
			"a request to save a valid project",
			"success",
			200,
			api.ProjectResp{"mock/project-x", "https://jabroni.com/mock/project-x", "jabroni.com"},
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
		{
			"",
			"a request to get provider projects",
			"a list of projects from the provider",
			200,
			[]provider.ProjectData{
				{"http://jabroni.com/mock/project1", "mock/project1"},
				{"http://jabroni.com/mock/project2", "mock/project2"},
				{"http://jabroni.com/mock/project3", "mock/project3"},
			},
			"mock_seeded_access_token_123",
		},
	}
	RunEndpointTests(t, "GET", "jabroni.com/projects", tests)

	invalidProviderTests := []*EndpointTest{
		afterAuthTest.UnsupportedProviderTest("", "jnk"),
	}
	RunEndpointTests(t, "GET", "jnk/projects", invalidProviderTests)
}

func TestGetProjects(t *testing.T) {
	tests := []*EndpointTest{
		MissingAuthorizationHeaderEndpointTest(""),
		InvalidAuthorizationTokenErrorTest(""),
		{
			"",
			"a request to get projects added to mspl",
			"a list of mspl projects",
			200,
			[]api.ProjectResp{
				{"mattmocks/project-1", "https://jabroni.com/mattmocks/project-1", "jabroni.com"},
				{"mattmocks/project-2", "https://jabroni.com/mattmocks/project-2", "jabroni.com"},
				{"mattmocks/project-3", "https://jabroni.com/mattmocks/project-3", "jabroni.com"},
			},
			"mock_seeded_access_token_123",
		},
	}
	RunEndpointTests(t, "GET", "projects", tests)
}

func TestDeleteProject(t *testing.T) {
	afterAuthTest := &AfterAuthTest{"mock_seeded_access_token_123"}
	tests := []*EndpointTest{
		MissingAuthorizationHeaderEndpointTest(""),
		InvalidAuthorizationTokenErrorTest(""),
		afterAuthTest.MissingParamEndpointTest(`{ "handle":"asdf" }`, "provider"),
		afterAuthTest.MissingParamEndpointTest(`{ "provider":"jabroni.com" }`, "handle"),
		afterAuthTest.UnsupportedProviderTest(`{ "handle":"mattmocks/project-1", "provider":"jnk" }`, "jnk"),
		{
			`{"provider":"jabroni.com", "handle":"mock/project-that-does-not-exist"}`,
			"a request to delete a project that does not exist",
			"a project not found response",
			400,
			api.NewNotFoundErr("project", "mock/project-that-does-not-exist"),
			"mock_seeded_access_token_123",
		},
	}
	RunEndpointTests(t, "DELETE", "projects", tests)

	et := &EndpointTest{
		`{"provider":"jabroni.com", "handle":"mattmocks/project-1"}`,
		"a request to delete a project added to mspl",
		"a success response",
		200,
		api.NewSuccessResponse(),
		"mock_seeded_access_token_123",
	}
	RunEndpointTests(t, "DELETE", "projects", []*EndpointTest{et})
	assertNumUserProjects(t, "After deleting a project", "mock_seeded_access_token_123", 2)
}
