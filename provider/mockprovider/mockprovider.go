package mockprovider

import (
	"io/ioutil"
	"net/http"

	"github.com/mcfly-svc/mcfly/provider"
	"github.com/stretchr/testify/mock"
)

type MockProvider struct {
	mock.Mock
}

func (p *MockProvider) Key() string {
	args := p.Called()
	return args.String(0)
}

func (p *MockProvider) GetProjects(token string, username string) ([]provider.ProjectData, error) {
	args := p.Called(token, username)
	return args.Get(0).([]provider.ProjectData), args.Error(1)
}

// get data from the provider based on a provider auth token
func (p *MockProvider) GetTokenData(token string) (*provider.TokenDataResponse, error) {
	args := p.Called(token)
	return args.Get(0).(*provider.TokenDataResponse), args.Error(1)
}

func (p *MockProvider) GetProjectData(token string, projectHandle string) (*provider.ProjectData, error) {
	args := p.Called(token, projectHandle)
	r0 := args.Get(0)
	r1 := args.Error(1)
	if r0 != nil {
		return r0.(*provider.ProjectData), r1
	} else {
		return nil, r1
	}
}

func (p *MockProvider) GetBuildData(token, buildHandle, projectHandle string) (*provider.BuildData, error) {
	args := p.Called(token, buildHandle, projectHandle)
	return args.Get(0).(*provider.BuildData), args.Error(1)
}

func (p *MockProvider) GetBuildConfig(token, buildHandle, projectHandle string) (*provider.BuildConfig, error) {
	args := p.Called(token, buildHandle, projectHandle)
	return args.Get(0).(*provider.BuildConfig), args.Error(1)
}

func (p *MockProvider) CreateProjectUpdateHook(token string, projectHandle string) error {
	args := p.Called(token, projectHandle)
	return args.Error(0)
}

func (p *MockProvider) DecodeProjectUpdateRequest(req *http.Request) (*provider.ProjectUpdateData, error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	args := p.Called(string(b))
	r0 := args.Get(0)
	if r0 != nil {
		return args.Get(0).(*provider.ProjectUpdateData), args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func strPtr(v string) *string { return &v }
