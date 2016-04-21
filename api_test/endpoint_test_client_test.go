package api_test

import (
  "net/http"
  "testing"
  "encoding/json"
  "github.com/mikec/marsupi-api/client"
)

type Entity struct {
  ID    int64   `json:"id"`
}

type EndpointTestClient struct {
	Test 				*testing.T
	Endpoint 		client.EntityEndpoint
}

func (self *EndpointTestClient) Create(JSON string) *http.Response {
	res, err := self.Endpoint.Create(JSON)
  if err != nil {
    self.Test.Error(err)
  }
  return res
}

func (self *EndpointTestClient) GetAllEntities() []Entity {
  res, err := self.Endpoint.GetAll()
  if err != nil {
    self.Test.Error(err)
  }
  var entities []Entity
  if err := json.NewDecoder(res.Body).Decode(&entities); err != nil {
    self.Test.Error(err)
  }
  return entities
}

func (self *EndpointTestClient) Get(ID int64) *http.Response {
	res, err := self.Endpoint.Get(ID)
  if err != nil {
    self.Test.Error(err)
  }
  return res
}

func (self *EndpointTestClient) GetEntity(ID int64) Entity {
  res := self.Get(ID)
  var e Entity
  if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
    self.Test.Error(err)
  }
  return e
}

func (self *EndpointTestClient) Delete(ID int64) *http.Response {
	res, err := self.Endpoint.Delete(ID)
	if err != nil {
		self.Test.Error(err)
	}
	return res
}
