package api_test

import (
	"testing"
	"fmt"
	"io/ioutil"
)

func TestLoginNewUserWithGitHubToken(t *testing.T) {
	cleanupDB()

	c := Client{t}
	_, res := c.Login("mock_github_token_123")

	b, err := ioutil.ReadAll(res.Body)
  if err != nil {
    t.Error(err)
  }

	fmt.Printf("LOGIN RESP:%+v\n",string(b))

}
