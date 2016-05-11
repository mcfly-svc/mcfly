package provider_test

import (
	"fmt"
	"testing"

	"github.com/google/go-github/github"
	"github.com/mikec/msplapi/provider"
	"github.com/stretchr/testify/assert"
)

var mockyJoe3Repos *github.RepositoriesSearchResult
var johnyNoRepos *github.RepositoriesSearchResult

func init() {
	mockyJoe3Repos = &github.RepositoriesSearchResult{
		intPtr(3),
		[]github.Repository{
			github.Repository{HTMLURL: strPtr("http://github.com/mockyjoe/mock.repo-1"), FullName: strPtr("mockyjoe/mock.repo-1")},
			github.Repository{HTMLURL: strPtr("http://github.com/mockyjoe/mock.repo-2"), FullName: strPtr("mockyjoe/mock.repo-2")},
			github.Repository{HTMLURL: strPtr("http://github.com/mockyjoe/mock.repo-3"), FullName: strPtr("mockyjoe/mock.repo-3")},
		},
	}
	johnyNoRepos = &github.RepositoriesSearchResult{
		intPtr(0),
		[]github.Repository{},
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
	if err := checkMockTokenInvalid(token); err != nil {
		return nil, nil, err
	}
	if repo == "does_not_exist" {
		return nil, nil, provider.NewProjectNotFoundErr("github", "mock/does_not_exist")
	}
	url := "http://github.com/mock/out"
	return &github.Repository{HTMLURL: &url}, nil, nil
}

func (self *MockGitHubClient) SearchRepos(
	token string,
	query string,
) (*github.RepositoriesSearchResult, *github.Response, error) {
	if err := checkMockTokenInvalid(token); err != nil {
		return nil, nil, err
	}
	if query == "user:johny_norepos" {
		return johnyNoRepos, nil, nil
	}
	return mockyJoe3Repos, nil, nil
}

func checkMockTokenInvalid(token string) error {
	if token == "mock_invalid_gh_token" {
		return provider.NewTokenInvalidErr("github")
	}
	return nil
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
	assertTokenInvalidErr(t, err)
}

func TestGetProjectDataNotFoundError(t *testing.T) {
	_, err := gh.GetProjectData("abc", "mock/does_not_exist")
	expectErrMsg := provider.NewProjectNotFoundErr("github", "mock/does_not_exist").Error()
	assert.Equal(t, expectErrMsg, err.Error())
}

func TestGetProjects(t *testing.T) {
	projects, _ := gh.GetProjects("mock_token", "mock_username")
	assert.Equal(t, 3, len(projects))
	assert.Equal(t, "http://github.com/mockyjoe/mock.repo-2", projects[1].Url)
	assert.Equal(t, "mockyjoe/mock.repo-2", projects[1].Handle)
}

func TestGetProjectsNoRepos(t *testing.T) {
	projects, _ := gh.GetProjects("mock_token", "johny_norepos")
	assert.Equal(t, 0, len(projects))
}

func TestGetProjectsInvalidToken(t *testing.T) {
	_, err := gh.GetProjects("mock_invalid_gh_token", "mockyjoe")
	assertTokenInvalidErr(t, err)
}

func assertTokenInvalidErr(t *testing.T, err error) {
	expectErrMsg := provider.NewTokenInvalidErr("github").Error()
	assert.Equal(t, expectErrMsg, err.Error())
}

func strPtr(v string) *string { return &v }
func intPtr(v int) *int       { return &v }
