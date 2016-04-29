package provider

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
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

type GitHub struct{}

func (self *GitHub) Key() string {
	return "github"
}

func (self *GitHub) GetTokenData(token string) (*TokenDataResponse, error) {
	gh := newClient(token)

	user, _, err := gh.Users.Get("")
	if err != nil {
		return nil, err
	}

	d := &TokenDataResponse{
		true,
		self.Key(),
		*user.Login,
		*user.Name,
	}

	return d, nil
}

func newClient(token string) *github.Client {
	tc := oauth2.NewClient(oauth2.NoContext, &TokenSource{token})
	return github.NewClient(tc)
}
