package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitHubClient interface {
	GetCurrentUser(string) (*github.User, *github.Response, error)
	GetRepo(string, string, string) (*github.Repository, *github.Response, error)
	SearchRepos(string, string) (*github.RepositoriesSearchResult, *github.Response, error)
	CreateHook(string, string, string, *github.Hook) (*github.Hook, *github.Response, error)
}

type GoGitHubClient struct{}

func (self *GoGitHubClient) GetCurrentUser(token string) (*github.User, *github.Response, error) {
	gh := self.NewClient(token)
	return gh.Users.Get("")
}

func (self *GoGitHubClient) GetRepo(
	token string,
	owner string,
	repo string,
) (*github.Repository, *github.Response, error) {
	gh := self.NewClient(token)
	return gh.Repositories.Get(owner, repo)
}

func (self *GoGitHubClient) SearchRepos(
	token string,
	query string,
) (*github.RepositoriesSearchResult, *github.Response, error) {
	gh := self.NewClient(token)
	return gh.Search.Repositories(query, &github.SearchOptions{})
}

func (self *GoGitHubClient) CreateHook(
	token,
	owner,
	repo string,
	hook *github.Hook,
) (*github.Hook, *github.Response, error) {
	gh := self.NewClient(token)
	return gh.Repositories.CreateHook(owner, repo, hook)
}

func (self *GoGitHubClient) NewClient(token string) *github.Client {
	tc := oauth2.NewClient(oauth2.NoContext, &TokenSource{token})
	return github.NewClient(tc)
}

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type GitHub struct {
	GitHubClient
	*SourceProviderConfig
}

func (self *GitHub) Key() string {
	return "github"
}

func (self *GitHub) GetTokenData(token string) (*TokenDataResponse, error) {
	user, _, err := self.GetCurrentUser(token)
	if err != nil {
		ghErr, ok := err.(*github.ErrorResponse)
		if !ok {
			return nil, err
		}
		if ghErr.Message == "Bad credentials" {
			return &TokenDataResponse{false, self.Key(), "", ""}, nil
		}
		return nil, ghErr
	}

	return &TokenDataResponse{true, self.Key(), *user.Login, *user.Name}, nil
}

func (self *GitHub) GetProjectData(token string, projectHandle string) (*ProjectData, error) {
	ph, err := NewProjectHandle(projectHandle)
	if err != nil {
		return nil, NewInvalidProjectHandleErr(self.Key(), projectHandle)
	}
	repo, _, err := self.GetRepo(token, ph.Owner, ph.Repo)
	if err != nil {
		return nil, self.handleGetProjectDataError(err, projectHandle)
	}
	return &ProjectData{*repo.HTMLURL, projectHandle}, nil
}

func (self *GitHub) GetProjects(token string, username string) ([]ProjectData, error) {
	repoSearchResult, _, err := self.SearchRepos(token, fmt.Sprintf("user:%s", username))
	if err != nil {
		return nil, self.handleGetProjectsError(err)
	}
	pData := make([]ProjectData, len(repoSearchResult.Repositories))
	for i, r := range repoSearchResult.Repositories {
		pData[i] = ProjectData{*r.HTMLURL, *r.FullName}
	}
	return pData, nil
}

func (self *GitHub) CreateProjectUpdateHook(token string, projectHandle string) error {
	ph, err := NewProjectHandle(projectHandle)
	if err != nil {
		return NewInvalidProjectHandleErr(self.Key(), projectHandle)
	}
	_, _, err = self.CreateHook(token, ph.Owner, ph.Repo, &github.Hook{
		Name:   strPtr("web"),
		Active: boolPtr(true),
		Events: []string{
			"push",
			"pull_request",
		},
		Config: map[string]interface{}{
			"url":          self.GetProjectUpdateHookUrl(self.Key()),
			"content_type": "json",
		},
	})
	if err != nil {
		// TODO: might need to handle the "Validation Failed" error that occurs
		// when the hook already exists
		return err
	}
	return nil
}

func (self *GitHub) DecodeProjectUpdateRequest(req *http.Request) (*ProjectUpdateData, error) {
	var payload github.WebHookPayload
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		return nil, err
	}
	pu := ProjectUpdateData{
		ProjectHandle: *payload.Repo.FullName,
		Builds:        make([]string, len(payload.Commits)),
	}
	for i, commit := range payload.Commits {
		pu.Builds[i] = *commit.ID
	}
	return &pu, nil
}

func (self *GitHub) handleGetProjectDataError(err error, projectHandle string) error {
	switch v := err.(type) {
	case *github.ErrorResponse:
		if ghErr := self.handleGitHubError(v); ghErr != nil {
			return ghErr
		}
		switch v.Message {
		case "Not Found":
			return NewProjectNotFoundErr(self.Key(), projectHandle)
		default:
			return v
		}
	default:
		return v
	}
}

func (self *GitHub) handleGetProjectsError(err error) error {
	switch v := err.(type) {
	case *github.ErrorResponse:
		if ghErr := self.handleGitHubError(v); ghErr != nil {
			return ghErr
		}
		switch v.Message {
		case "Validation Failed":
			return NewGetProjectsFailedErr(self.Key())
		default:
			return v
		}
	default:
		return v
	}
}

func (self *GitHub) handleGitHubError(ghErr *github.ErrorResponse) error {
	switch ghErr.Message {
	case "Bad credentials":
		return NewTokenInvalidErr(self.Key())
	default:
		return nil
	}
}

type ProjectHandle struct {
	Owner string
	Repo  string
}

func NewProjectHandle(projectHandle string) (*ProjectHandle, error) {
	s := strings.Split(projectHandle, "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("Invalid project handle `%s`", projectHandle)
	}
	return &ProjectHandle{s[0], s[1]}, nil
}

func (ph ProjectHandle) String() string {
	return fmt.Sprintf("%s/%s", ph.Owner, ph.Repo)
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
