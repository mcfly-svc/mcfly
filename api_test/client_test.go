package api_test

import (
  "testing"
  "net/http"
  "github.com/mikec/marsupi-api/client"
)

type Client struct {
	Test 				*testing.T
}

func (self *Client) Login(token string) (*client.ClientResponse, *http.Response) {
	cr, res, err := apiClient.Login(token)
  if err != nil {
    self.Test.Error(err)
  }
  return cr, res
}
