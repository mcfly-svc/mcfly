package api_test

import (
  "net/http"
  "testing"
  "github.com/mikec/marsupi-api/client"
)

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

func (self *EndpointTestClient) GetAll() []interface{} {
	res, err := self.Endpoint.GetAll()
  if err != nil {
    self.Test.Error(err)
  }
  entities, err := self.Endpoint.Decoder.DecodeArrayResponse(res)
  if err != nil {
  	self.Test.Error(err)
  }
  return entities
}

func (self *EndpointTestClient) GetRes(ID int64) *http.Response {
	res, err := self.Endpoint.Get(ID)
  if err != nil {
    self.Test.Error(err)
  }
  return res
}

func (self *EndpointTestClient) Get(ID int64) interface{} {
	res := self.GetRes(ID)
  entity, err := self.Endpoint.Decoder.DecodeResponse(res)
  if err != nil {
  	self.Test.Error(err)
  }
  return entity
}

func (self *EndpointTestClient) Delete(ID int64) *http.Response {
	res, err := self.Endpoint.Delete(ID)
	if err != nil {
		self.Test.Error(err)
	}
	return res
}
