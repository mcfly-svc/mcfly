package client

import (
	"github.com/mikec/msplapi/api/apidata"
	"github.com/mikec/msplapi/provider"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type Client interface {
	Login(*apidata.LoginReq) (*ClientResponse, *http.Response, error)
	AddProject(*apidata.ProjectReq) (*ClientResponse, *http.Response, error)
	DeleteProject(*apidata.ProjectReq) (*ClientResponse, *http.Response, error)
	Deploy(*apidata.DeployReq) (*ClientResponse, *http.Response, error)
	SaveBuild(*apidata.BuildReq) (*ClientResponse, *http.Response, error)
	GetProviderProjects(string) (*ClientResponse, *http.Response, error)
	GetProjects() (*ClientResponse, *http.Response, error)
}

type ClientResponse struct {
	Data       interface{}
	StatusCode int
}

type MsplClient struct {
	ServerURL   string
	AccessToken string
}

func NewMsplClient(serverURL string, token string) *MsplClient {
	return &MsplClient{serverURL, token}
}

func (self *MsplClient) Login(request *apidata.LoginReq) (*ClientResponse, *http.Response, error) {
	return self.DoMsplClientRequest("POST", "login", request, &apidata.LoginResp{})
}

func (self *MsplClient) AddProject(request *apidata.ProjectReq) (*ClientResponse, *http.Response, error) {
	return self.DoMsplClientRequest("POST", "projects", request, &apidata.ApiResponse{})
}

func (self *MsplClient) DeleteProject(request *apidata.ProjectReq) (*ClientResponse, *http.Response, error) {
	return self.DoMsplClientRequest("DELETE", "projects", request, &apidata.ApiResponse{})
}

func (self *MsplClient) Deploy(request *apidata.DeployReq) (*ClientResponse, *http.Response, error) {
	return self.DoMsplClientRequest("POST", "deploy", request, &apidata.ApiResponse{})
}

func (self *MsplClient) SaveBuild(request *apidata.BuildReq) (*ClientResponse, *http.Response, error) {
	return self.DoMsplClientRequest("POST", "builds", request, &apidata.ApiResponse{})
}

func (self *MsplClient) GetProviderProjects(providerKey string) (*ClientResponse, *http.Response, error) {
	res, err := self.DoGet(fmt.Sprintf("%s/projects", providerKey), nil)
	if err != nil {
		return nil, nil, err
	}

	var pData []provider.ProjectData
	return decodeResponse(res, &pData)
}

func (self *MsplClient) GetProjects() (*ClientResponse, *http.Response, error) {
	res, err := self.DoGet("projects", nil)
	if err != nil {
		return nil, nil, err
	}

	var projects []apidata.ProjectResp
	return decodeResponse(res, &projects)
}

func (self *MsplClient) DoMsplClientRequest(
	method, endpoint string,
	reqData interface{},
	respData interface{},
) (*ClientResponse, *http.Response, error) {
	json, err := json.Marshal(reqData)
	if err != nil {
		return nil, nil, err
	}
	jsonStr := string(json)
	res, err := self.DoReq(method, endpoint, &jsonStr, nil)
	if err != nil {
		return nil, nil, err
	}
	return decodeResponse(res, respData)
}

func (self *MsplClient) EndpointUrl(endpointName string) string {
	return fmt.Sprintf("%s/api/0/%s", self.ServerURL, endpointName)
}

func (self *MsplClient) DoGet(endpoint string, opts *RequestOptions) (*http.Response, error) {
	return self.DoReq("GET", endpoint, nil, opts)
}

func (self *MsplClient) DoPost(endpoint string, JSON *string, opts *RequestOptions) (*http.Response, error) {
	return self.DoReq("POST", endpoint, JSON, opts)
}

func (self *MsplClient) DoDelete(endpoint string, JSON *string, opts *RequestOptions) (*http.Response, error) {
	return self.DoReq("DELETE", endpoint, JSON, opts)
}

func (self *MsplClient) DoReq(method string, endpoint string, JSON *string, opts *RequestOptions) (*http.Response, error) {
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

func decodeResponse(res *http.Response, v interface{}) (*ClientResponse, *http.Response, error) {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res, err
	}
	res.Body = bodyReader{bytes.NewBuffer(b)}
	if err := json.Unmarshal(b, &v); err != nil {
		apiErr := apidata.ApiError{}
		if uerr := json.Unmarshal(b, &apiErr); uerr != nil {
			return nil, res, NewBodyDecodeError(string(b))
		}
		return &ClientResponse{apiErr, res.StatusCode}, res, nil
	}
	return &ClientResponse{v, res.StatusCode}, res, nil
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
		var apiErr apidata.ApiError
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
