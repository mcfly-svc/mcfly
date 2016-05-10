package provider

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GitHubClient interface {
	GetCurrentUser(string) (*github.User, *github.Response, error)
	GetRepo(string, string, string) (*github.Repository, *github.Response, error)
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
	ph, err := parseProjectHandle(projectHandle)
	if err != nil {
		return nil, NewProjectDataInvalidHandleErr(projectHandle, self.Key())
	}
	repo, _, err := self.GetRepo(token, ph.Owner, ph.Repo)
	if err != nil {
		ghErr, ok := err.(*github.ErrorResponse)
		if !ok {
			return nil, err
		}
		switch ghErr.Message {
		case "Not Found":
			return nil, NewProjectDataNotFoundErr(projectHandle, self.Key())
		case "Bad credentials":
			return nil, NewProjectDataTokenInvalidErr(projectHandle, self.Key())
		default:
			return nil, ghErr
		}
	}
	return &ProjectData{*repo.HTMLURL, projectHandle}, nil
}

func (self *GitHub) GetProjects(token string) ([]ProjectData, error) {
	// TMP
	return []ProjectData{
		{"http://jnk1", "gh/jnk1"},
		{"http://jnk2", "gh/jnk2"},
		{"http://jnk3", "gh/jnk3"},
	}, nil
}

type ProjectHandle struct {
	Owner string
	Repo  string
}

func parseProjectHandle(projectHandle string) (*ProjectHandle, error) {
	s := strings.Split(projectHandle, "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("Invalid project handle `%s`", projectHandle)
	}
	return &ProjectHandle{s[0], s[1]}, nil
}
