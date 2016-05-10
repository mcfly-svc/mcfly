package api_test

import (
	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/client"
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
	dbUrl     string
)

type MockLogger struct{}

func (m MockLogger) Handler(h http.Handler, s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

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
		return nil, provider.NewProjectDataNotFoundErr("mock/project-x", "jabroni.com")
	}
	if projectHandle == "invalid_project_handle" {
		return nil, provider.NewProjectDataInvalidHandleErr("invalid_project_handle", "jabroni.com")
	}
	if token == "bad_saved_jabroni.com_token_123" {
		return nil, provider.NewProviderTokenInvalidErr("jabroni.com")
	}
	return &provider.ProjectData{"https://jabroni.com/mock/project-x", "mock/project-x"}, nil
}

func init() {

	jabroni := MockProvider{}

	dbUrl = "postgres://localhost:5432/marsupi_test?sslmode=disable"

	authProviders := make(map[string]provider.AuthProvider)
	authProviders[jabroni.Key()] = &jabroni

	sourceProviders := make(map[string]provider.SourceProvider)
	sourceProviders[jabroni.Key()] = &jabroni

	server = httptest.NewServer(
		api.NewRouter(
			dbUrl,
			//MockLogger{},
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

func resetDB() {
	cleanupDB()
	seedDB()
}

func cleanupDB() {
	db.Clean(dbUrl)
}

func seedDB() {
	db.Seed(dbUrl)
}
