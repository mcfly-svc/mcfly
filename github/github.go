package github

import (
  "golang.org/x/oauth2"
  "github.com/google/go-github/github"
  "github.com/mikec/marsupi-api/models"
)

type TokenSource struct {
  AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
  token := &oauth2.Token{
    AccessToken: t.AccessToken,
  }
  return token, nil
}

type Client struct {}

func (c Client) GetAuthenticatedUser(token string) (*models.User, error) {
  gh := newClient(token)

  user, _, err := gh.Users.Get("")
  if err != nil {
    return nil, err
  }

  u := &models.User{
    Name: *user.Name,
    GitHubUsername: *user.Login,
    GitHubToken: token,
  }
  return u, nil
}

func newClient(token string) *github.Client {
  tc := oauth2.NewClient(oauth2.NoContext, &TokenSource{token})
  return github.NewClient(tc)
}