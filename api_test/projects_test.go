package api_test

import "testing"

func TestPostProject(t *testing.T) {

	validJson := `{ "project_name":"jabroni.com/mockers/mock-project", "provider":"jabroni.com" }`

	tests := []*EndpointTest{
		InvalidJsonEndpointTest(),

		MissingAuthorizationHeaderEndpointTest(&validJson),

		InvalidAuthorizationTokenErrorTest(&validJson),

		MissingParamEndpointTest(`{ "project_name":"asdf" }`, "provider"),

		MissingParamEndpointTest(`{ "provider":"jabroni.com" }`, "project_name"),
	}

	RunEndpointTests(t, "POST", "projects", tests)

}
