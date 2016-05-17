package api_test

import (
	"fmt"
	"log"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/client"
	"github.com/mikec/msplapi/config"
	"github.com/mikec/msplapi/db"
	"github.com/mikec/msplapi/logging"
	"github.com/mikec/msplapi/models"
	"github.com/mikec/msplapi/mq"
	"github.com/mikec/msplapi/provider"
	"github.com/streadway/amqp"

	"io"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	server    *httptest.Server
	reader    io.Reader
	apiClient *client.Client
	cfg       *config.Config
	database  *models.DB
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

type MockMessageChannel struct{}

func (m *MockMessageChannel) Send(q *amqp.Queue, v interface{}) error {
	return nil
}

func (m *MockMessageChannel) SendDeployQueueMessage(dpq *mq.DeployQueueMessage) error {
	if dpq.BuildHandle == "start_deploy_error" {
		return fmt.Errorf("mock start deploy error")
	}
	return nil
}

func (m *MockMessageChannel) CloseConnection() error {
	return nil
}

func (m *MockMessageChannel) CloseChannel() error {
	return nil
}

func init() {
	cfg = GetTestConfig()
	dbs, err := models.NewDB(cfg.DatabaseUrl)

	database = dbs
	if err != nil {
		log.Fatal(err)
	}
	createDB()

	jabroni := MockProvider{}

	msgChannel := &MockMessageChannel{}

	authProviders := make(map[string]provider.AuthProvider)
	authProviders[jabroni.Key()] = &jabroni

	sourceProviders := make(map[string]provider.SourceProvider)
	sourceProviders[jabroni.Key()] = &jabroni

	server = httptest.NewServer(
		api.NewRouter(
			cfg,
			logging.HttpRequestLogger{},
			msgChannel,
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

func createDB() {
	db.Create(cfg.DatabaseUrl)
}

func resetDB() {
	cleanupDB()
	seedDB()
}

func cleanupDB() {
	db.Clean(database.DB)
}

func seedDB() {
	db.Seed(database.DB)
}
