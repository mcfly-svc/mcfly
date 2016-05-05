package provider_test

import (
	"testing"

	"github.com/google/go-github/github"
	"github.com/mikec/msplapi/provider"
	"github.com/stretchr/testify/assert"
)

type MockGitHubClient struct{}

func (self *MockGitHubClient) GetCurrentUser(token string) (*github.User, *github.Response, error) {
	if token == "mock_bad_token" {
		return nil, nil, &github.ErrorResponse{Message: "Bad credentials"}
	} else {
		return &github.User{Login: strPtr("@mjones"), Name: strPtr("Mock Jones")}, nil, nil
	}
}

var gh provider.GitHub

func init() {
	gh = provider.GitHub{&MockGitHubClient{}}
}

func TestGetTokenDataBadCredentialsError(t *testing.T) {
	td, _ := gh.GetTokenData("mock_bad_token")
	assert.False(t, td.IsValid, "Token data should not be valid")
}

func TestGetTokenDataValidToken(t *testing.T) {
	td, _ := gh.GetTokenData("mock_good_token")
	assert.True(t, td.IsValid, "Token data should be valid")
	assert.Equal(t, "github", td.Provider, "Token data response should include provider name")
	assert.Equal(t, "@mjones", td.ProviderUsername, "Token data response should include provider username")
	assert.Equal(t, "Mock Jones", td.UserName, "Token data response should include user's name")
}

func strPtr(s string) *string { return &s }
