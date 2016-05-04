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

// get data from the provider based on a provider auth token
func (ap *MockAuthProvider) GetTokenData(token string) (*provider.TokenDataResponse, error) {
	if token == "badtoken" {
		return &provider.TokenDataResponse{false, ap.Key(), "", ""}, nil
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

func cleanupDB() {
	db.Clean(dbUrl)
}

func seedDB() {
	db.Seed(dbUrl)
}
