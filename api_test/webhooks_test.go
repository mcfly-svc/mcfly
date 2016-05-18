package api_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/provider"
	"github.com/stretchr/testify/assert"
)

func (p *MockProvider) CreateProjectUpdateHook(token string, projectHandle string) error {
	return nil
}

func (p *MockProvider) DecodeProjectUpdateRequest(req *http.Request) (*provider.ProjectUpdateData, error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	switch string(b) {
	case "valid_two_commits":
		return &provider.ProjectUpdateData{
			ProjectHandle: "mattmocks/project-3",
			Builds:        []string{"abc", "123"},
		}, nil
	case "project_handle_dne":
		return &provider.ProjectUpdateData{
			ProjectHandle: "jnk/project-dne",
			Builds:        []string{},
		}, nil
	case "decode_error":
		return nil, fmt.Errorf("mock decode error")
	default:
		return nil, nil
	}
}

func TestProjectUpdateWebhookErrors(t *testing.T) {
	cleanupDB()
	RunEndpointTests(t, "POST", "webhooks/jabroni.com/project-update", []*EndpointTest{
		{
			`decode_error`,
			"a request with an invalid payload that causes decoding to fail",
			"a server error",
			400,
			api.NewServerErr(),
			"",
		},
		{
			`project_handle_dne`,
			"a request with a valid payload for a project that does not exist",
			"a server error",
			400,
			api.NewServerErr(),
			"",
		},
	})
}

func TestProjectUpdateWebhook(t *testing.T) {
	cleanupDB()
	RunEndpointTests(t, "POST", "webhooks/jabroni.com/project-update", []*EndpointTest{
		{
			`valid_two_commits`,
			"a request with a valid payload with two commits",
			"success",
			200,
			api.NewSuccessResponse(),
			"",
		},
	})

	assert.Equal(t, "mattmocks/project-3", _lastDeployQueueMessage.ProjectHandle,
		"should send the correct project handle in the deploy queue message")

	assert.Equal(t, "jabroni.com", _lastDeployQueueMessage.SourceProvider,
		"should send the correct provider in the deploy queue message")

	assert.Equal(t, "123", _lastDeployQueueMessage.BuildHandle,
		"should send the correct build handle in the deploy queue message")

	assert.Equal(t, "mock_jabroni.com_token_123", _lastDeployQueueMessage.SourceProviderAccessToken,
		"should send the correct provider access token in the deploy queue message")

}
