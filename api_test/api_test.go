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
	apiClient client.Client
	dbUrl     string
)

type MockLogger struct{}

func (m MockLogger) Handler(h http.Handler, s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

type MockAuthProvider struct{}

func (ap *MockAuthProvider) Key() string {
	return "jabroni.com"
}

func generateMockToken() string {
	return "mock_generated_access_token_123"
}

// get data from the provider based on a provider auth token
func (ap *MockAuthProvider) GetTokenData(token string) (*provider.TokenDataResponse, error) {
	if token == "badtoken" {
		return &provider.TokenDataResponse{false, ap.Key(), "", ""}, nil
	} else if token == "mock_jabroni.com_token_123" {
		return &provider.TokenDataResponse{true, ap.Key(), "mattmocks", "Matt Mockman"}, nil
	}
	return &provider.TokenDataResponse{true, ap.Key(), "mikej", "Mike Jimmers"}, nil
}

func init() {

	jabroni := MockAuthProvider{}

	dbUrl = "postgres://localhost:5432/marsupi_test?sslmode=disable"

	authProviders := make(map[string]provider.AuthProvider)
	authProviders[jabroni.Key()] = &jabroni

	server = httptest.NewServer(
		api.NewRouter(
			dbUrl,
			//MockLogger{},
			logging.HttpRequestLogger{},
			generateMockToken,
			authProviders,
		),
	)
	apiClient = client.NewClient(server.URL)
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
