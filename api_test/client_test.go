package api_test

import (
  "bytes"
  "net/http"
  "io"
  "io/ioutil"
  "testing"
  "encoding/json"
  "github.com/mikec/marsupi-api/client"
)

type Entity struct {
  ID    int64   `json:"id"`
}

type Client struct {
	Test 				*testing.T
	Endpoint 		client.EntityEndpoint
}

type bodyReader struct {
    *bytes.Buffer
}

func (m bodyReader) Close() error { return nil } 

func (self *Client) Create(JSON string) (Entity, *http.Response) {
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

func (self *Client) GetAll() ([]Entity, *http.Response) {
  res, err := self.Endpoint.GetAll()
  if err != nil {
    self.Test.Error(err)
  }
  var entities []Entity
  if err := decodeEntitiesResponse(res, &entities); err != nil {
    self.Test.Error(err)
  }
  return entities, res
}

func (self *Client) Get(ID int64) (Entity, *http.Response) {
	res, err := self.Endpoint.Get(ID)
  if err != nil {
    self.Test.Error(err)
  }
  var e Entity
  if err := decodeEntityResponse(res, &e); err != nil {
    self.Test.Error(err)
  }
  return e, res
}

func (self *Client) Delete(ID int64) *http.Response {
	res, err := self.Endpoint.Delete(ID)
	if err != nil {
		self.Test.Error(err)
	}
	return res
}

func decodeEntitiesResponse(res *http.Response, entities *[]Entity) error {
  err := decodeBody(res, entities)
  if err != nil {
    return err
  }
  return nil
}

func decodeEntityResponse(res *http.Response, entity *Entity) error {
  err := decodeBody(res, entity)
  if err != nil {
    return err
  }
  return nil
}

func decodeBody(res *http.Response, v interface{}) error {
  body, err := readBody(res)
  if err != nil {
    return err
  }
  if err := json.NewDecoder(body).Decode(v); err != nil {
    return err
  }
  return nil
}

func readBody(res *http.Response) (io.ReadCloser, error) {
  buf, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }
  body := bodyReader{bytes.NewBuffer(buf)}
  res.Body = bodyReader{bytes.NewBuffer(buf)}

  return body, nil
}
