package api_test

import (
	"github.com/mikec/marsupi-api/api"
	"github.com/mikec/marsupi-api/client"
	"github.com/mikec/marsupi-api/db"
	"github.com/mikec/marsupi-api/logging"
	"github.com/mikec/marsupi-api/provider"

	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	server    *httptest.Server
	reader    io.Reader
	apiClient client.Client
)

type MockLogger struct{}

func (m MockLogger) Handler(h http.Handler, s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

type MockAuthProvider struct{}

func (ap MockAuthProvider) Key() string {
	return "jabroni.com"
}

// get data from the provider based on a provider auth token
func (ap MockAuthProvider) GetTokenData(string) (*provider.TokenDataResponse, error) {
	return &provider.TokenDataResponse{true, ap.Key(), "mikej", "Mike Jimmers"}, nil
}

func init() {

	jabroni := MockAuthProvider{}

	authProviders := make(map[string]provider.AuthProvider)
	authProviders[jabroni.Key()] = &jabroni

	server = httptest.NewServer(
		api.NewRouter(
			"postgres://localhost:5432/marsupi_test?sslmode=disable",
			//MockLogger{},
			logging.HttpRequestLogger{},
			authProviders,
		),
	)
	apiClient = client.NewClient(server.URL)
}

func TestMain(m *testing.M) {
	log.Println("Recreating the database")
	recreateDB()

	ret := m.Run()

	log.Println("Cleaning up the database")
	cleanupDB()

	os.Exit(ret)
}

func recreateDB() {
	db.RunHelperScript("../db/helpers/recreate.sh")
}

func cleanupDB() {
	db.RunHelperScript("../db/helpers/clean.sh")
}
