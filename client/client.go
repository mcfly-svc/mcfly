package client

import (
  "io"
	"net/http"
	"fmt"
	"strings"
)

type EntityEndpoint struct {
  Name                      string
  ServerURL                 string
  Decoder                   EntityDecoder
}

type Client struct {
  Projects      EntityEndpoint
  Users         EntityEndpoint
}

func NewClient(serverURL string) Client {
  return Client{
    EntityEndpoint{"projects", serverURL, ProjectDecoder{}},
    EntityEndpoint{"users", serverURL, ProjectDecoder{}},
  }
}

func (self *EntityEndpoint) EndpointUrl() string {
  return fmt.Sprintf("%s/api/0/%s", self.ServerURL, self.Name)
}

func (self *EntityEndpoint) EntityUrl(id int64) string {
  return fmt.Sprintf("%s/%d", self.EndpointUrl(), id)
}

func (self *EntityEndpoint) Create(JSON string) (*http.Response, error) {
  res, err := DoPost(self.EndpointUrl(), JSON)
  if err != nil {
    return nil, err
  }
  return res, nil
}

func (self *EntityEndpoint) Delete(id int64) (*http.Response, error) {
  res, err := DoDelete(self.EntityUrl(id))
  if err != nil {
    return nil, err
  }
  return res, nil
}

func (self *EntityEndpoint) GetAll() (*http.Response, error) {
  res, err := DoGet(self.EndpointUrl())
  if err != nil {
    return nil, err
  }
  return res, nil
}

func (self *EntityEndpoint) Get(id int64) (*http.Response, error) {
  res, err := DoGet(self.EntityUrl(id))
  if err != nil {
    return nil, err
  }
  return res, nil

}

func DoGet(url string) (*http.Response, error) {
  return DoReq("GET", url, nil)
}

func DoPost(url string, JSON string) (*http.Response, error) {
  return DoReq("POST", url, &JSON)
}

func DoDelete(url string) (*http.Response, error) {
  return DoReq("DELETE", url, nil)
}

func DoReq(method string, url string, JSON *string) (*http.Response, error) {
  var reader io.Reader
  if JSON != nil {
    reader = strings.NewReader(*JSON)
  }

  req, err := http.NewRequest(method, url, reader)
  if err != nil {
    return nil, err
  }

  res, err := http.DefaultClient.Do(req)
  if err != nil {
    return nil, err
  }

  return res, nil
}

