package client

import (
  "github.com/mikec/marsupi-api/api"
  "github.com/mikec/marsupi-api/models"

  "bytes"
  "io"
  "io/ioutil"
  "encoding/json"
	"net/http"
	"fmt"
	"strings"
  "reflect"
)

type ClientResponse struct {
  Data            interface{}
  StatusCode      int
}

type EntityEndpoint struct {
  Name                      string
  SingularName              string
  ServerURL                 string
  ModelType                 reflect.Type
}

type Client struct {
  Projects      EntityEndpoint
  Users         EntityEndpoint
}

func NewClient(serverURL string) Client {
  return Client{
    EntityEndpoint{"projects", "project", serverURL, reflect.TypeOf(models.Project{})},
    EntityEndpoint{"users", "user", serverURL, reflect.TypeOf(models.User{})},
  }
}

func (self *EntityEndpoint) EndpointUrl() string {
  return fmt.Sprintf("%s/api/0/%s", self.ServerURL, self.Name)
}

func (self *EntityEndpoint) EntityUrl(id int64) string {
  return fmt.Sprintf("%s/%d", self.EndpointUrl(), id)
}

func (self *EntityEndpoint) Create(JSON string) (*ClientResponse, *http.Response, error) {
  res, err := DoPost(self.EndpointUrl(), JSON)
  if err != nil {
    return nil, nil, err
  }

  d, err := decodeResponse(res, self.ModelType)
  if err != nil {
    return nil, res, err
  }
  cr := &ClientResponse{d, res.StatusCode}

  return cr, res, nil
}

func (self *EntityEndpoint) Delete(id int64) (*ClientResponse, *http.Response, error) {
  res, err := DoDelete(self.EntityUrl(id))
  if err != nil {
    return nil, nil, err
  }

  d, err := decodeResponse(res, reflect.TypeOf(api.ApiResponse{}))
  if err != nil {
    return nil, res, err
  }
  return &ClientResponse{d, res.StatusCode}, res, nil
}

func (self *EntityEndpoint) GetAll() (*ClientResponse, *http.Response, error) {
  res, err := DoGet(self.EndpointUrl())
  if err != nil {
    return nil, nil, err
  }

  d, err := decodeResponse(res, reflect.SliceOf(self.ModelType))
  if err != nil {
    return nil, res, err
  }
  return &ClientResponse{d, res.StatusCode}, res, nil
}

func (self *EntityEndpoint) Get(id int64) (*ClientResponse, *http.Response, error) {
  res, err := DoGet(self.EntityUrl(id))
  if err != nil {
    return nil, nil, err
  }
  d, err := decodeResponse(res, self.ModelType)
  if err != nil {
    return nil, res, err
  }
  return &ClientResponse{d, res.StatusCode}, res, nil
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

type bodyReader struct {
    *bytes.Buffer
}

func (m bodyReader) Close() error { return nil } 

func decodeResponse(res *http.Response, t reflect.Type) (interface{}, error) {
  b, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }
  res.Body = bodyReader{bytes.NewBuffer(b)}
  var v interface{}
  json.Unmarshal(b, &v)
  m, ok := v.(map[string]interface{})
  if ok && m["error"] != nil {
    var apiErr api.ApiError
    json.Unmarshal(b, &apiErr)
    return apiErr, nil
  } else {
    m := reflect.New(t).Interface()
    json.Unmarshal(b, m)
    return m, nil
  }
}
