package provider

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	ListHooks(string, string, string) ([]github.Hook, *github.Response, error)
	DeleteHook(string, string, string, int) (*github.Response, error)
	GetCommit(string, string, string, string) (*github.RepositoryCommit, *github.Response, error)
}

type GoGitHubClient struct{}

func (self *GoGitHubClient) GetCurrentUser(token string) (*github.User, *github.Response, error) {
	gh := self.NewClient(token)
	return gh.Users.Get("")
}

func (self *GoGitHubClient) GetRepo(
	token, owner, repo string,
) (*github.Repository, *github.Response, error) {
	gh := self.NewClient(token)
	return gh.Repositories.Get(owner, repo)
}

func (self *GoGitHubClient) SearchRepos(
	token, query string,
) (*github.RepositoriesSearchResult, *github.Response, error) {
	gh := self.NewClient(token)
	return gh.Search.Repositories(query, &github.SearchOptions{})
}

func (self *GoGitHubClient) CreateHook(
	token, owner, repo string, hook *github.Hook,
) (*github.Hook, *github.Response, error) {
	gh := self.NewClient(token)
	return gh.Repositories.CreateHook(owner, repo, hook)
}

func (self *GoGitHubClient) DeleteHook(
	token, owner, repo string, id int,
) (*github.Response, error) {
	gh := self.NewClient(token)
	return gh.Repositories.DeleteHook(owner, repo, id)
}

func (self *GoGitHubClient) ListHooks(token, owner, repo string) ([]github.Hook, *github.Response, error) {
	gh := self.NewClient(token)
	return gh.Repositories.ListHooks(owner, repo, nil)
}

func (self *GoGitHubClient) GetCommit(token, owner, repo, sha string,
) (*github.RepositoryCommit, *github.Response, error) {
	gh := self.NewClient(token)
	return gh.Repositories.GetCommit(owner, repo, sha)
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
			return &TokenDataResponse{false, self.Key(), "", nil}, nil
		}
		return nil, ghErr
	}
	return &TokenDataResponse{true, self.Key(), *user.Login, user.Name}, nil
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

func (self *GitHub) GetBuildData(token, sha, projectHandle string) (*BuildData, error) {
	ph, err := NewProjectHandle(projectHandle)
	if err != nil {
		return nil, NewInvalidProjectHandleErr(self.Key(), projectHandle)
	}
	commit, _, err := self.GetCommit(token, ph.Owner, ph.Repo, sha)
	if err != nil {
		return nil, err
	}
	return &BuildData{
		Url:    ptrToStr(commit.HTMLURL),
		Handle: sha,
		Config: []byte(""),
	}, nil
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

	existingHook, err := self.GetProjectUpdateHook(token, ph.Owner, ph.Repo)
	if err != nil {
		return err
	}
	if existingHook != nil {
		_, err = self.DeleteHook(token, ph.Owner, ph.Repo, *existingHook.ID)
		if err != nil {
			return err
		}
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
			"secret":       self.SourceProviderConfig.WebhookSecret,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (self *GitHub) DecodeProjectUpdateRequest(req *http.Request) (*ProjectUpdateData, error) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = bodyReader{bytes.NewBuffer(bodyBytes)}

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

	xHubSig := req.Header.Get("X-Hub-Signature")
	whSecret := self.SourceProviderConfig.WebhookSecret

	signatureVerified := checkSig(bodyBytes, xHubSig, whSecret)
	if !signatureVerified {
		return nil, NewInvalidWebhookSignatureErr(self.Key())
	}

	return &pu, nil
}

func (self *GitHub) GetProjectUpdateHook(token, owner, repo string) (*github.Hook, error) {
	hooks, _, err := self.ListHooks(token, owner, repo)
	if err != nil {
		return nil, err
	}
	hookUrl := self.GetProjectUpdateHookUrl(self.Key())
	for _, hook := range hooks {
		if hook.Config["url"] == hookUrl {
			return &hook, nil
		}
	}
	return nil, nil
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

func checkSig(body []byte, signature, key string) bool {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write(body)
	expectedSig := fmt.Sprintf("sha1=%s", hex.EncodeToString(mac.Sum(nil)))
	return signature == expectedSig
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

type bodyReader struct {
	*bytes.Buffer
}

func (m bodyReader) Close() error { return nil }

func ptrToStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
