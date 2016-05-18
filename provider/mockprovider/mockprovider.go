package mockprovider

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mikec/msplapi/provider"
)

type MockProvider struct{}

func (ap *MockProvider) Key() string {
	return "jabroni.com"
}

func (self *MockProvider) GetProjects(token string, username string) ([]provider.ProjectData, error) {
	return []provider.ProjectData{
		{"http://jabroni.com/mock/project1", "mock/project1"},
		{"http://jabroni.com/mock/project2", "mock/project2"},
		{"http://jabroni.com/mock/project3", "mock/project3"},
	}, nil
}

// get data from the provider based on a provider auth token
func (p *MockProvider) GetTokenData(token string) (*provider.TokenDataResponse, error) {
	if token == "badtoken" {
		return &provider.TokenDataResponse{false, p.Key(), "", nil}, nil
	} else if token == "mock_jabroni.com_token_123" {
		return &provider.TokenDataResponse{true, p.Key(), "mattmocks", strPtr("Matt Mockman")}, nil
	}
	return &provider.TokenDataResponse{true, p.Key(), "mikej", strPtr("Mike Jimmers")}, nil
}

func (p *MockProvider) GetProjectData(token string, projectHandle string) (*provider.ProjectData, error) {
	if projectHandle == "project_handle_dne" {
		return nil, provider.NewProjectNotFoundErr("jabroni.com", "mock/project-x")
	}
	if projectHandle == "invalid_project_handle" {
		return nil, provider.NewInvalidProjectHandleErr("jabroni.com", "invalid_project_handle")
	}
	if token == "bad_saved_jabroni.com_token_123" {
		return nil, provider.NewTokenInvalidErr("jabroni.com")
	}
	return &provider.ProjectData{"https://jabroni.com/mock/project-x", "mock/project-x"}, nil
}

func (self *MockProvider) GetBuildData(token, buildHandle, projectHandle string) (*provider.BuildData, error) {
	return nil, nil
}

func (p *MockProvider) CreateProjectUpdateHook(token string, projectHandle string) error {
	return nil
}

func (p *MockProvider) DecodeProjectUpdateRequest(req *http.Request) (*provider.ProjectUpdateData, error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	switch string(b) {
	case "valid_two_commits":
		return &provider.ProjectUpdateData{
			ProjectHandle: "mattmocks/project-3",
			Builds:        []string{"abc", "123"},
		}, nil
	case "project_handle_dne":
		return &provider.ProjectUpdateData{
			ProjectHandle: "jnk/project-dne",
			Builds:        []string{},
		}, nil
	case "decode_error":
		return nil, fmt.Errorf("mock decode error")
	default:
		return nil, nil
	}
}

func strPtr(v string) *string { return &v }
