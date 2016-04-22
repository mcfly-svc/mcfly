package api_test

import (
  "testing"
  "github.com/mikec/marsupi-api/client"
)

type Client struct {
	Test 				*testing.T
	Endpoint 		client.EntityEndpoint
}

func (self *Client) Create(JSON string) (*client.ClientResponse) {
	res, err := self.Endpoint.Create(JSON)
  if err != nil {
    self.Test.Error(err)
  }
  return res
}

func (self *Client) GetAll() (*client.ClientResponse) {
  res, err := self.Endpoint.GetAll()
  if err != nil {
    self.Test.Error(err)
  }
  return res
}

func (self *Client) Get(ID int64) (*client.ClientResponse) {
	res, err := self.Endpoint.Get(ID)
  if err != nil {
    self.Test.Error(err)
  }
  return res
}

func (self *Client) Delete(ID int64) (*client.ClientResponse) {
	res, err := self.Endpoint.Delete(ID)
	if err != nil {
		self.Test.Error(err)
	}
	return res
}
