package api_test

import (
	"testing"

	"github.com/mikec/msplapi/api"
)

func TestDeploy(t *testing.T) {

	afterAuthTest := &AfterAuthTest{"mock_seeded_access_token_123"}
	validJson := `{ "build_handle": "a1b2c3", "project_handle":"mattmocks/project-2", "provider":"jabroni.com" }`
	startDeployErr := `{ "build_handle": "start_deploy_error", "project_handle":"mattmocks/project-2", "provider":"jabroni.com" }`
	unsupportedProvider := `{ "build_handle": "a1b2c3", "project_handle":"mattmocks/project-2", "provider":"jnk" }`
	missingBuildHandle := `{ "project_handle":"mattmocks/project-2", "provider":"jabroni.com" }`
	missingProjectHandle := `{ "build_handle": "a1b2c3", "provider":"jabroni.com" }`
	missingProvider := `{ "build_handle": "a1b2c3", "project_handle":"mattmocks/project-2" }`
	projectDne := `{ "build_handle": "a1b2c3", "project_handle":"jnk/jnk", "provider":"jabroni.com" }`

	tests := []*EndpointTest{
		afterAuthTest.InvalidJsonEndpointTest(),
		MissingAuthorizationHeaderEndpointTest(validJson),
		InvalidAuthorizationTokenErrorTest(validJson),
		afterAuthTest.MissingParamEndpointTest(missingBuildHandle, "build_handle"),
		afterAuthTest.MissingParamEndpointTest(missingProjectHandle, "project_handle"),
		afterAuthTest.MissingParamEndpointTest(missingProvider, "provider"),
		afterAuthTest.UnsupportedProviderTest(unsupportedProvider, "jnk"),
		{
			startDeployErr,
			"a request to deploy when SendStartDeployMessage returns an error",
			"a server error",
			400,
			api.NewServerErr(),
			"mock_seeded_access_token_123",
		},
		{
			projectDne,
			"a request to deploy for a project that does not exist",
			"a project not found error",
			400,
			api.NewNotFoundErr("project", "jnk/jnk"),
			"mock_seeded_access_token_123",
		},
		{
			validJson,
			"a request to deploy",
			"a success response",
			200,
			api.NewSuccessResponse(),
			"mock_seeded_access_token_123",
		},
	}

	RunEndpointTests(t, "POST", "deploy", tests)

}
