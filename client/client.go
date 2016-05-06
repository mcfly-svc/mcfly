package client

import (
	"github.com/mikec/msplapi/api"

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
	ServerURL   string
	AccessToken string
}

type ClientResponse struct {
	Data       interface{}
	StatusCode int
}

func NewClient(serverURL string, token string) *Client {
	return &Client{serverURL, token}
}

func (self *Client) Login(token string, provider string) (*ClientResponse, *http.Response, error) {
	json, err := json.Marshal(api.LoginReq{token, provider})
	if err != nil {
		return nil, nil, err
	}
	jsonStr := string(json)
	res, err := self.DoPost("login", &jsonStr, nil)
	if err != nil {
		return nil, nil, err
	}

	var loginResp api.LoginResp
	if err := decodeResponse(res, &loginResp); err != nil {
		return nil, res, err
	}

	return &ClientResponse{loginResp, res.StatusCode}, res, nil
}

func (self *Client) AddProject(projectHandle string, provider string) (*ClientResponse, *http.Response, error) {
	json, err := json.Marshal(api.PostProjectReq{projectHandle, provider})
	if err != nil {
		return nil, nil, err
	}
	jsonStr := string(json)
	res, err := self.DoPost("projects", &jsonStr, nil)
	if err != nil {
		return nil, nil, err
	}

	var apiResp api.ApiResponse
	if err := decodeResponse(res, &apiResp); err != nil {
		return nil, res, err
	}

	return &ClientResponse{apiResp, res.StatusCode}, res, nil
}

func (self *Client) EndpointUrl(endpointName string) string {
	return fmt.Sprintf("%s/api/0/%s", self.ServerURL, endpointName)
}

func (self *Client) DoGet(endpoint string, opts *RequestOptions) (*http.Response, error) {
	return self.DoReq("GET", endpoint, nil, opts)
}

func (self *Client) DoPost(endpoint string, JSON *string, opts *RequestOptions) (*http.Response, error) {
	return self.DoReq("POST", endpoint, JSON, opts)
}

func (self *Client) DoDelete(endpoint string, opts *RequestOptions) (*http.Response, error) {
	return self.DoReq("DELETE", endpoint, nil, opts)
}

func (self *Client) DoReq(method string, endpoint string, JSON *string, opts *RequestOptions) (*http.Response, error) {
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

	if self.AccessToken != "" && opts.AuthHeader == nil {
		opts.AuthHeader = &self.AccessToken
	}

	if opts.AuthHeader != nil {
		req.Header.Add("Authorization", *opts.AuthHeader)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type RequestOptions struct {
	UseBasicAuth bool
	AuthHeader   *string
}

func NewRequestOptions() *RequestOptions {
	return &RequestOptions{true, nil}
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
