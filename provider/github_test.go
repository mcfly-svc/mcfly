package provider_test

import (
	"fmt"
	"testing"

	"github.com/google/go-github/github"
	"github.com/mikec/msplapi/provider"
	"github.com/stretchr/testify/assert"
)

var mockyJoe3Repos *github.RepositoriesSearchResult

func init() {
	mockyJoe3Repos = &github.RepositoriesSearchResult{
		intPtr(3),
		[]github.Repository{
			github.Repository{HTMLURL: strPtr("http://github.com/mockyjoe/mock.repo-1"), FullName: strPtr("mockyjoe/mock.repo-1")},
			github.Repository{HTMLURL: strPtr("http://github.com/mockyjoe/mock.repo-2"), FullName: strPtr("mockyjoe/mock.repo-2")},
			github.Repository{HTMLURL: strPtr("http://github.com/mockyjoe/mock.repo-3"), FullName: strPtr("mockyjoe/mock.repo-3")},
		},
	}
}

type MockGitHubClient struct{}

func (self *MockGitHubClient) GetCurrentUser(token string) (*github.User, *github.Response, error) {
	if token == "mock_bad_token" {
		return nil, nil, &github.ErrorResponse{Message: "Bad credentials"}
	} else {
		return &github.User{Login: strPtr("@mjones"), Name: strPtr("Mock Jones")}, nil, nil
	}
}

func (self *MockGitHubClient) GetRepo(
	token string,
	owner string,
	repo string,
) (*github.Repository, *github.Response, error) {
	if token == "mock_invalid_gh_token" {
		return nil, nil, provider.NewProviderTokenInvalidErr("github")
	}
	if repo == "does_not_exist" {
		return nil, nil, provider.NewProjectDataNotFoundErr("mock/does_not_exist", "github")
	}
	url := "http://github.com/mock/out"
	return &github.Repository{HTMLURL: &url}, nil, nil
}

func (self *MockGitHubClient) GetReposByOwner(
	token string,
	owner string,
) (*github.RepositoriesSearchResult, *github.Response, error) {
	if token == "mock_invalid_gh_token" {
		return nil, nil, provider.NewProviderTokenInvalidErr("github")
	}
	return mockyJoe3Repos, nil, nil
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

func TestGetProjectData(t *testing.T) {

	var tests = []struct {
		Handle                    string
		ExpectValidHandle         bool
		ExpectProjectDataReturned bool
	}{
		{"has/oneslash", true, true},
		{"noslash", false, false},
		{"has/two/slashes", false, false},
	}

	for _, tst := range tests {
		p, err := gh.GetProjectData("abc", tst.Handle)
		if tst.ExpectValidHandle {
			assert.Nil(t, err, fmt.Sprintf("Expected `%s` to be a valid project handle", tst.Handle))
		} else {
			assert.NotNil(t, err, fmt.Sprintf("Expected `%s` to be an invalid project handle", tst.Handle))
		}
		if tst.ExpectProjectDataReturned {
			assert.Equal(t, "http://github.com/mock/out", p.Url)
		}
	}
}

func TestGetProjectDataInvalidTokenError(t *testing.T) {
	_, err := gh.GetProjectData("mock_invalid_gh_token", "mock/project")
	expectErrMsg := provider.NewProviderTokenInvalidErr("github").Error()
	assert.Equal(t, expectErrMsg, err.Error())
}

func TestGetProjectDataNotFoundError(t *testing.T) {
	_, err := gh.GetProjectData("abc", "mock/does_not_exist")
	expectErrMsg := provider.NewProjectDataNotFoundErr("mock/does_not_exist", "github").Error()
	assert.Equal(t, expectErrMsg, err.Error())
}

func TestGetProjects(t *testing.T) {

}

func TestGetProjectsInvalidToken(t *testing.T) {
	_, err := gh.GetProjects("mock_invalid_gh_token", "mockyjoe")
	expectErrMsg := provider.NewProviderTokenInvalidErr("github").Error()
	assert.Equal(t, expectErrMsg, err.Error())
}

func strPtr(v string) *string { return &v }
func intPtr(v int) *int       { return &v }
