package mockclient

import (
	"fmt"
	"net/http"

	"github.com/mikec/msplapi/api/apidata"
	"github.com/mikec/msplapi/client"
)

type MockClient struct{}

func (c *MockClient) Login(request *apidata.LoginReq) (*client.ClientResponse, *http.Response, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func (c *MockClient) AddProject(request *apidata.ProjectReq) (*client.ClientResponse, *http.Response, error) {
	return nil, nil, fmt.Errorf("not implemented")
}
func (c *MockClient) DeleteProject(request *apidata.ProjectReq) (*client.ClientResponse, *http.Response, error) {
	return nil, nil, fmt.Errorf("not implemented")
}
func (c *MockClient) Deploy(request *apidata.DeployReq) (*client.ClientResponse, *http.Response, error) {
	return nil, nil, fmt.Errorf("not implemented")
}
func (c *MockClient) SaveBuild(request *apidata.BuildReq) (*client.ClientResponse, *http.Response, error) {
	return nil, nil, fmt.Errorf("not implemented")
}
func (c *MockClient) GetProviderProjects(string) (*client.ClientResponse, *http.Response, error) {
	return nil, nil, fmt.Errorf("not implemented")
}
func (c *MockClient) GetProjects() (*client.ClientResponse, *http.Response, error) {
	return nil, nil, fmt.Errorf("not implemented")
}
