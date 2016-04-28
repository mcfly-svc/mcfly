package api_test

import (
	"github.com/mikec/marsupi-api/client"
	"net/http"
	"testing"
)

type Client struct {
	Test *testing.T
}

func (self *Client) Login(token string, tokenType string) (*client.ClientResponse, *http.Response) {
	cr, res, err := apiClient.Login(token, tokenType)
	if err != nil {
		self.Test.Error(err)
	}
	return cr, res
}
