package client

import (
	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/models"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type Client struct {
	Context  *ClientContext
	Projects EntityEndpoint
	Users    EntityEndpoint
}

type ClientResponse struct {
	Data       interface{}
	StatusCode int
}

func NewClient(serverURL string) Client {
	ctx := ClientContext{serverURL}
	return Client{
		&ctx,
		EntityEndpoint{"projects", "project", &ctx, reflect.TypeOf(models.Project{})},
		EntityEndpoint{"users", "user", &ctx, reflect.TypeOf(models.User{})},
	}
}

func (self *Client) Login(token string, provider string) (*ClientResponse, *http.Response, error) {
	json, err := json.Marshal(api.LoginReq{token, provider})
	if err != nil {
		return nil, nil, err
	}
	jsonStr := string(json)
	res, err := self.Context.DoPost("login", &jsonStr, nil)
	if err != nil {
		return nil, nil, err
	}

	var loginResp api.LoginResp
	if err := decodeResponse(res, &loginResp); err != nil {
		return nil, res, err
	}

	return &ClientResponse{loginResp, res.StatusCode}, res, nil
}

type EntityEndpoint struct {
	Name         string
	SingularName string
	Context      *ClientContext
	ModelType    reflect.Type
}

func (self *EntityEndpoint) Create(JSON string) (*ClientResponse, *http.Response, error) {
	res, err := self.Context.DoPost(self.Name, &JSON, nil)
	if err != nil {
		return nil, nil, err
	}

	d, err := decodeResponseFromType(res, self.ModelType)
	if err != nil {
		return nil, res, err
	}
	cr := &ClientResponse{d, res.StatusCode}

	return cr, res, nil
}

func (self *EntityEndpoint) Delete(id int64) (*ClientResponse, *http.Response, error) {
	res, err := self.Context.DoDelete(fmt.Sprintf("%s/%d", self.Name, id), nil)
	if err != nil {
		return nil, nil, err
	}

	d, err := decodeResponseFromType(res, reflect.TypeOf(api.ApiResponse{}))
	if err != nil {
		return nil, res, err
	}
	return &ClientResponse{d, res.StatusCode}, res, nil
}

func (self *EntityEndpoint) GetAll() (*ClientResponse, *http.Response, error) {
	res, err := self.Context.DoGet(self.Name, nil)
	if err != nil {
		return nil, nil, err
	}

	d, err := decodeResponseFromType(res, reflect.SliceOf(self.ModelType))
	if err != nil {
		return nil, res, err
	}
	return &ClientResponse{d, res.StatusCode}, res, nil
}

func (self *EntityEndpoint) Get(id int64) (*ClientResponse, *http.Response, error) {
	res, err := self.Context.DoGet(fmt.Sprintf("%s/%d", self.Name, id), nil)
	if err != nil {
		return nil, nil, err
	}
	d, err := decodeResponseFromType(res, self.ModelType)
	if err != nil {
		return nil, res, err
	}
	return &ClientResponse{d, res.StatusCode}, res, nil
}

type RequestOptions struct {
	UseBasicAuth bool
}

func NewRequestOptions() *RequestOptions {
	return &RequestOptions{true}
}

type ClientContext struct {
	ServerURL string
}

func (self *ClientContext) EndpointUrl(endpointName string) string {
	return fmt.Sprintf("%s/api/0/%s", self.ServerURL, endpointName)
}

func (self *ClientContext) DoGet(endpoint string, opts *RequestOptions) (*http.Response, error) {
	return self.DoReq("GET", endpoint, nil, opts)
}

func (self *ClientContext) DoPost(endpoint string, JSON *string, opts *RequestOptions) (*http.Response, error) {
	return self.DoReq("POST", endpoint, JSON, opts)
}

func (self *ClientContext) DoDelete(endpoint string, opts *RequestOptions) (*http.Response, error) {
	return self.DoReq("DELETE", endpoint, nil, opts)
}

func (self *ClientContext) DoReq(method string, endpoint string, JSON *string, opts *RequestOptions) (*http.Response, error) {
	if opts == nil {
		opts = NewRequestOptions()
	}

	url := self.EndpointUrl(endpoint)

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

func decodeResponse(res *http.Response, v interface{}) error {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body = bodyReader{bytes.NewBuffer(b)}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	return nil
}

func decodeResponseFromType(res *http.Response, t reflect.Type) (interface{}, error) {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body = bodyReader{bytes.NewBuffer(b)}
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	m, ok := v.(map[string]interface{})
	if ok && m["error"] != nil {
		var apiErr api.ApiError
		if err := json.Unmarshal(b, &apiErr); err != nil {
			return nil, err
		}
		return apiErr, nil
	} else {
		m := reflect.New(t).Interface()
		if err := json.Unmarshal(b, m); err != nil {
			return nil, err
		}
		return m, nil
	}
}
