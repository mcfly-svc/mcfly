package mockclient

import (
	"errors"
	"net/http"

	"github.com/mikec/msplapi/api/apidata"
	"github.com/mikec/msplapi/client"
	"github.com/stretchr/testify/mock"
)

var ErrMock = errors.New("mock error")
var ClientResponseError = &client.ClientResponse{
	StatusCode: 400,
	Data:       &apidata.ApiError{Error: "mock api error"},
}
var ClientResponseSuccess = &client.ClientResponse{
	StatusCode: 200,
	Data:       &apidata.ApiResponse{"success!"},
}

type MockClient struct {
	mock.Mock
}

func (c *MockClient) Login(request *apidata.LoginReq) (*client.ClientResponse, *http.Response, error) {
	args := c.Called(request)
	return args.Get(0).(*client.ClientResponse), args.Get(1).(*http.Response), args.Error(2)
}
func (c *MockClient) AddProject(request *apidata.ProjectReq) (*client.ClientResponse, *http.Response, error) {
	args := c.Called(request)
	return args.Get(0).(*client.ClientResponse), args.Get(1).(*http.Response), args.Error(2)
}
func (c *MockClient) DeleteProject(request *apidata.ProjectReq) (*client.ClientResponse, *http.Response, error) {
	args := c.Called(request)
	return args.Get(0).(*client.ClientResponse), args.Get(1).(*http.Response), args.Error(2)
}
func (c *MockClient) Deploy(request *apidata.DeployReq) (*client.ClientResponse, *http.Response, error) {
	args := c.Called(request)
	return args.Get(0).(*client.ClientResponse), args.Get(1).(*http.Response), args.Error(2)
}
func (c *MockClient) SaveBuild(request *apidata.BuildReq) (*client.ClientResponse, *http.Response, error) {
	args := c.Called(request)
	return args.Get(0).(*client.ClientResponse), args.Get(1).(*http.Response), args.Error(2)
}
func (c *MockClient) GetProviderProjects(providerKey string) (*client.ClientResponse, *http.Response, error) {
	args := c.Called(providerKey)
	return args.Get(0).(*client.ClientResponse), args.Get(1).(*http.Response), args.Error(2)
}
func (c *MockClient) GetProjects() (*client.ClientResponse, *http.Response, error) {
	args := c.Called()
	return args.Get(0).(*client.ClientResponse), args.Get(1).(*http.Response), args.Error(2)
}
