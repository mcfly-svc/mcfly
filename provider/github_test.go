package provider_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/github"
	"github.com/mcfly-svc/mcfly/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type MockGitHubClient struct {
	mock.Mock
}

func (self *MockGitHubClient) GetCurrentUser(token string) (*github.User, *github.Response, error) {
	if token == "mock_bad_token" {
		return nil, nil, &github.ErrorResponse{Message: "Bad credentials"}
	} else if token == "nil_github_name" {
		return &github.User{Login: strPtr("@mjones"), Name: nil}, nil, nil
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

func (self *MockGitHubClient) CreateHook(
	token,
	owner,
	repo string,
	hook *github.Hook,
) (*github.Hook, *github.Response, error) {
	args := self.Called(token, owner, repo, hook)
	r0, _ := args.Get(0).(*github.Hook)
	r1, _ := args.Get(1).(*github.Response)
	return r0, r1, args.Error(2)
}

func (self *MockGitHubClient) DeleteHook(
	token,
	owner,
	repo string,
	id int,
) (*github.Response, error) {
	args := self.Called(token, owner, repo, id)
	r0, _ := args.Get(0).(*github.Response)
	return r0, args.Error(1)
}

func (self *MockGitHubClient) ListHooks(token, owner, repo string) ([]github.Hook, *github.Response, error) {
	args := self.Called(token, owner, repo)
	r0, _ := args.Get(0).([]github.Hook)
	r1, _ := args.Get(1).(*github.Response)
	return r0, r1, args.Error(2)
}

func (self *MockGitHubClient) GetCommit(token, owner, repo, sha string,
) (*github.RepositoryCommit, *github.Response, error) {
	return nil, nil, nil
}

func (self *MockGitHubClient) GetContents(token, owner, repo, path, sha string) (
	*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
	return nil, nil, nil, nil
}

func checkMockTokenInvalid(token string) error {
	if token == "mock_invalid_gh_token" {
		return provider.NewTokenInvalidErr("github")
	}
	return nil
}

var gh provider.GitHub

func init() {
	gh = newGhProvider(nil)
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
	assert.Equal(t, strPtr("Mock Jones"), td.UserName, "Token data response should include user's name")
}

func TestGetTokenDataNilGitHubName(t *testing.T) {
	td, _ := gh.GetTokenData("nil_github_name")
	assert.True(t, td.IsValid, "Token data should be valid")
	assert.Equal(t, "github", td.Provider, "Token data response should include provider name")
	assert.Equal(t, "@mjones", td.ProviderUsername, "Token data response should include provider username")
	assert.Nil(t, td.UserName, "Token data response should include a nil UserName")
}

func TestGetProjectData(t *testing.T) {
	p, _ := gh.GetProjectData("abc", "has/oneslash")
	assert.Equal(t, "http://github.com/mock/out", p.Url)
	runProjectHandleValidationTests(t, func(handle string) error {
		_, err := gh.GetProjectData("abc", handle)
		return err
	})
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

func TestCreateProjectUpdateHook(t *testing.T) {
	tests := []struct {
		ExistingHookReturn    *github.Hook
		ExistingHookReturnErr error
		DeleteHookReturnErr   error
		CreateHookReturnErr   error
		ExpectReturn          error
	}{
		{
			ExistingHookReturn: &github.Hook{
				ID: intPtr(123),
				Config: map[string]interface{}{
					"url": "http://mocky.com/api/0/webhooks/github/project-update",
				},
			},
			ExistingHookReturnErr: nil,
			DeleteHookReturnErr:   nil,
			CreateHookReturnErr:   nil,
			ExpectReturn:          nil,
		},
		{
			ExistingHookReturn:    nil,
			ExistingHookReturnErr: errors.New("mock existing hook error"),
			DeleteHookReturnErr:   nil,
			CreateHookReturnErr:   nil,
			ExpectReturn:          errors.New("mock existing hook error"),
		},
		{
			ExistingHookReturn: &github.Hook{
				ID: intPtr(123),
				Config: map[string]interface{}{
					"url": "http://mocky.com/api/0/webhooks/github/project-update",
				},
			},
			ExistingHookReturnErr: nil,
			DeleteHookReturnErr:   errors.New("mock delete hook error"),
			CreateHookReturnErr:   nil,
			ExpectReturn:          errors.New("mock delete hook error"),
		},
		{
			ExistingHookReturn: &github.Hook{
				ID: intPtr(123),
				Config: map[string]interface{}{
					"url": "http://mocky.com/api/0/webhooks/github/project-update",
				},
			},
			ExistingHookReturnErr: nil,
			DeleteHookReturnErr:   nil,
			CreateHookReturnErr:   errors.New("mock create hook error"),
			ExpectReturn:          errors.New("mock create hook error"),
		},
	}

	for _, test := range tests {
		m := new(MockGitHubClient)
		if test.ExistingHookReturn != nil {
			m.On("ListHooks", "abc", "jim", "proj1").Return([]github.Hook{*test.ExistingHookReturn}, nil, nil)
		} else {
			m.On("ListHooks", "abc", "jim", "proj1").Return(nil, nil, test.ExistingHookReturnErr)
		}
		if test.ExistingHookReturn != nil {
			if test.DeleteHookReturnErr != nil {
				m.On("DeleteHook", "abc", "jim", "proj1", *test.ExistingHookReturn.ID).Return(
					nil, test.DeleteHookReturnErr)
			} else {
				m.On("DeleteHook", "abc", "jim", "proj1", *test.ExistingHookReturn.ID).Return(nil, nil)
			}
			if test.CreateHookReturnErr != nil {
				m.On("CreateHook", "abc", "jim", "proj1", mock.AnythingOfType("*github.Hook")).Return(
					nil, nil, test.CreateHookReturnErr)
			} else {
				m.On("CreateHook", "abc", "jim", "proj1", mock.AnythingOfType("*github.Hook")).Return(nil, nil, nil)
			}
		}
		ghp := newGhProvider(m)
		ret := ghp.CreateProjectUpdateHook("abc", "jim/proj1")
		assert.Equal(t, test.ExpectReturn, ret)
	}

}

func TestDecodeProjectUpdateRequest(t *testing.T) {
	tests := []struct {
		BodyString string
		Expect     *provider.ProjectUpdateData
		ExpectErr  bool
		Message    string
		ValidSig   bool
	}{
		{
			BodyString: "{junk}",
			Expect:     nil,
			ExpectErr:  true,
			Message:    "Should fail on invalid JSON",
			ValidSig:   true,
		},
		{
			BodyString: `{"repository":{"full_name":"mockman/banana"},"commits":[{"id":"abc"},{"id":"123"}]}`,
			Expect: &provider.ProjectUpdateData{
				Builds:        []string{"abc", "123"},
				ProjectHandle: "mockman/banana",
			},
			ExpectErr: false,
			Message:   "When given a payload with a two commits it should return the correct ProjectUpdateData",
			ValidSig:  true,
		},
	}
	for _, test := range tests {
		req := &http.Request{
			Header: make(map[string][]string),
			Body:   bodyReader{bytes.NewBuffer([]byte(test.BodyString))},
		}
		var secret string
		if test.ValidSig {
			secret = gh.SourceProviderConfig.WebhookSecret
		} else {
			secret = "invalid_secret"
		}
		req.Header.Add("X-Hub-Signature", makeSig(test.BodyString, secret))
		pData, err := gh.DecodeProjectUpdateRequest(req)
		assert.Equal(t, test.Expect, pData, test.Message)
		if test.ExpectErr {
			assert.NotNil(t, err, test.Message)
		}
	}
}

/*

// TODO: test GetBuildData and GetBuildConfig
//       and use testify mocks for MockGitHubClient

func TestGetBuildConfig(t *testing.T) {
	tests := []struct{
		Token string,
		Sha string,
		ProjectHandle string,
		ExpectBuildConfig *provider.BuildConfig,
		ExpectErr error
	}{
		{
			Token: "valid_token",
			Sha: "valid_sha",
			ProjectHandle:
		}
	}
}*/

func TestNewGitHubProjectHandle(t *testing.T) {
	tests := []struct {
		Handle              string
		ExpectProjectHandle *provider.GitHubProjectHandle
		ExpectErr           error
	}{
		{
			Handle: "valid/handle",
			ExpectProjectHandle: &provider.GitHubProjectHandle{
				Owner: "valid",
				Repo:  "handle",
			},
			ExpectErr: nil,
		},
		{
			Handle:              "invalid/handle/now",
			ExpectProjectHandle: nil,
			ExpectErr:           fmt.Errorf("Invalid project handle `invalid/handle/now`"),
		},
	}

	for _, test := range tests {
		ph, err := provider.NewGitHubProjectHandle(test.Handle)
		assert.Equal(t, test.ExpectErr, err)
		assert.Equal(t, test.ExpectProjectHandle, ph)
	}
}

func assertTokenInvalidErr(t *testing.T, err error) {
	expectErrMsg := provider.NewTokenInvalidErr("github").Error()
	assert.Equal(t, expectErrMsg, err.Error())
}

func runProjectHandleValidationTests(t *testing.T, fn func(string) error) {
	var tests = []struct {
		Handle            string
		ExpectValidHandle bool
	}{
		{"has/oneslash", true},
		{"noslash", false},
		{"has/two/slashes", false},
	}
	for _, tst := range tests {
		err := fn(tst.Handle)
		if tst.ExpectValidHandle {
			assert.Nil(t, err, fmt.Sprintf("Expected `%s` to be a valid project handle", tst.Handle))
		} else {
			assert.NotNil(t, err, fmt.Sprintf("Expected `%s` to be an invalid project handle", tst.Handle))
		}
	}
}

func newGhProvider(mockGhClient *MockGitHubClient) provider.GitHub {
	srcProviderCfg := provider.SourceProviderConfig{
		ProjectUpdateHookUrlFmt: fmt.Sprintf("http://mocky.com/api/0/webhooks/{provider}/project-update"),
		WebhookSecret:           "mock_webhook_secret",
	}
	if mockGhClient == nil {
		mockGhClient = new(MockGitHubClient)
	}
	return provider.GitHub{
		GitHubClient:         mockGhClient,
		SourceProviderConfig: &srcProviderCfg,
	}
}

type bodyReader struct {
	*bytes.Buffer
}

func (m bodyReader) Close() error { return nil }

func makeSig(message, key string) string {
	//fmt.Printf("ENCODING body/key %s/%s\n", message, key)
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(message))
	return fmt.Sprintf("sha1=%s", hex.EncodeToString(mac.Sum(nil)))
}
