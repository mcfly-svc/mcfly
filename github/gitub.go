package github

import (
  "golang.org/x/oauth2"
  "github.com/google/go-github/github"
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

func GetAuthClient(token string) (client *github.Client) {
	tc := oauth2.NewClient(oauth2.NoContext, &TokenSource{token})
  client = github.NewClient(tc)
  return
}
