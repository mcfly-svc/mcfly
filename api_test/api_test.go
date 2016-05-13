package api_test

import (
	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/client"
	"github.com/mikec/msplapi/config"
	"github.com/mikec/msplapi/db"
	"github.com/mikec/msplapi/logging"
	"github.com/mikec/msplapi/provider"

	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	server    *httptest.Server
	reader    io.Reader
	apiClient *client.Client
	cfg       *config.Config
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

func generateMockToken() string {
	return "mock_generated_access_token_123"
}

// get data from the provider based on a provider auth token
func (p *MockProvider) GetTokenData(token string) (*provider.TokenDataResponse, error) {
	if token == "badtoken" {
		return &provider.TokenDataResponse{false, p.Key(), "", ""}, nil
	} else if token == "mock_jabroni.com_token_123" {
		return &provider.TokenDataResponse{true, p.Key(), "mattmocks", "Matt Mockman"}, nil
	}
	return &provider.TokenDataResponse{true, p.Key(), "mikej", "Mike Jimmers"}, nil
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

func (p *MockProvider) CreateProjectUpdateHook(token string, projectHandle string) error {
	return nil
}

func (p *MockProvider) DecodeProjectUpdateRequest(req *http.Request) (*provider.ProjectUpdateData, error) {
	return nil, nil
}

func init() {
	cfg = GetTestConfig()
	resetDB()

	jabroni := MockProvider{}

	authProviders := make(map[string]provider.AuthProvider)
	authProviders[jabroni.Key()] = &jabroni

	sourceProviders := make(map[string]provider.SourceProvider)
	sourceProviders[jabroni.Key()] = &jabroni

	server = httptest.NewServer(
		api.NewRouter(
			cfg,
			logging.HttpRequestLogger{},
			generateMockToken,
			authProviders,
			sourceProviders,
		),
	)
	apiClient = client.NewClient(server.URL, "")
}

func TestMain(m *testing.M) {

	ret := m.Run()

	// setup

	os.Exit(ret)
}

func NewApiClient(token string) *client.Client {
	return client.NewClient(server.URL, token)
}

func resetDB() {
	cleanupDB()
	seedDB()
}

func cleanupDB() {
	db.Clean(cfg.DatabaseUrl)
}

func seedDB() {
	db.Seed(cfg.DatabaseUrl)
}
