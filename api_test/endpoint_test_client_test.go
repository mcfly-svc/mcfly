package api_test

import (
  "bytes"
  "net/http"
  "io/ioutil"
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

type bodyReader struct {
    *bytes.Buffer
}

func (m bodyReader) Close() error { return nil } 

func (self *EndpointTestClient) Create(JSON string) (Entity, *http.Response) {
	res, err := self.Endpoint.Create(JSON)
  if err != nil {
    self.Test.Error(err)
  }

  var e Entity
  if err := decodeEntityResponse(res, &e); err != nil {
    self.Test.Error(err)
  }

  return e, res
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

func decodeEntityResponse(res *http.Response, entity *Entity) error {
  buf, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return err
  }
  body := bodyReader{bytes.NewBuffer(buf)}

  if err := json.NewDecoder(body).Decode(entity); err != nil {
    return err
  }

  res.Body = bodyReader{bytes.NewBuffer(buf)}

  return nil
}
