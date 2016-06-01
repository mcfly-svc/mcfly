package api_test

import (
	"fmt"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mcfly-svc/mcfly/api"
	"github.com/mcfly-svc/mcfly/client"
	"github.com/mcfly-svc/mcfly/config"
	"github.com/mcfly-svc/mcfly/db"
	"github.com/mcfly-svc/mcfly/logging"
	"github.com/mcfly-svc/mcfly/mq"
	"github.com/mcfly-svc/mcfly/provider"
	"github.com/mcfly-svc/mcfly/provider/mockprovider"
	"github.com/streadway/amqp"

	"io"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	server    *httptest.Server
	reader    io.Reader
	apiClient *client.McflyClient
	cfg       *config.Config
	mdb       *db.McflyDB
	jabroni   *mockprovider.MockProvider
)

func generateMockToken() string {
	return "mock_generated_access_token_123"
}

type MockMessageChannel struct{}

func (m *MockMessageChannel) Send(q *amqp.Queue, v interface{}) error {
	return nil
}

var _lastDeployQueueMessage *mq.DeployQueueMessage

func (m *MockMessageChannel) SendDeployQueueMessage(dpq *mq.DeployQueueMessage) error {
	if dpq.BuildHandle == "start_deploy_error" {
		return fmt.Errorf("mock start deploy error")
	}
	_lastDeployQueueMessage = dpq
	return nil
}

func (m *MockMessageChannel) CloseConnection() error {
	return nil
}

func (m *MockMessageChannel) CloseChannel() error {
	return nil
}

func init() {
	fmt.Println("INIT")

	cfg = GetTestConfig()
	mcflyDb := db.NewMcflyDB(cfg.DatabaseUrl, cfg.DatabaseName, cfg.DatabaseUseSSL)
	mdb = mcflyDb
	initDB()

	jabroni = new(mockprovider.MockProvider)
	jabroni.On("Key").Return("jabroni.com")
	jabroni.On("GetProjects", "mock_jabroni.com_token_123", "mattmocks").Return([]provider.ProjectData{
		{"http://jabroni.com/mock/project1", "mock/project1"},
		{"http://jabroni.com/mock/project2", "mock/project2"},
		{"http://jabroni.com/mock/project3", "mock/project3"},
	}, nil)
	jabroni.On("GetTokenData", "badtoken").Return(
		&provider.TokenDataResponse{false, "jabroni.com", "", nil}, nil,
	)
	jabroni.On("GetTokenData", "mock_jabroni.com_token_123").Return(
		&provider.TokenDataResponse{true, "jabroni.com", "mattmocks", strPtr("Matt Mockman")}, nil,
	)
	jabroni.On("GetTokenData", "mock_dne_user_token_123").Return(
		&provider.TokenDataResponse{true, "jabroni.com", "mikej", strPtr("Mike Jimmers")}, nil,
	)
	jabroni.On("GetProjectData", "mock_jabroni.com_token_123", "project_handle_dne").Return(
		nil, provider.NewProjectNotFoundErr("jabroni.com", "mock/project-x"),
	)
	jabroni.On("GetProjectData", "mock_jabroni.com_token_123", "invalid_project_handle").Return(
		nil, provider.NewInvalidProjectHandleErr("jabroni.com", "invalid_project_handle"),
	)
	jabroni.On("GetProjectData", "bad_saved_jabroni.com_token_123", "mock/project-x").Return(
		nil, provider.NewTokenInvalidErr("jabroni.com"),
	)
	jabroni.On("GetProjectData", "mock_jabroni.com_token_123", "mattmocks/project-1").Return(
		&provider.ProjectData{"https://jabroni.com/mock/project-x", "mock/project-x"}, nil,
	)
	jabroni.On("GetProjectData", "mock_jabroni.com_token_123", "mock/project-x").Return(
		&provider.ProjectData{"https://jabroni.com/mock/project-x", "mock/project-x"}, nil,
	)
	jabroni.On("CreateProjectUpdateHook", "mock_jabroni.com_token_123", "mock/project-x").Return(nil)

	msgChannel := &MockMessageChannel{}

	authProviders := make(map[string]provider.AuthProvider)
	authProviders[jabroni.Key()] = jabroni

	sourceProviders := make(map[string]provider.SourceProvider)
	sourceProviders[jabroni.Key()] = jabroni

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
	apiClient = client.NewMcflyClient(server.URL, "")
}

func TestMain(m *testing.M) {
	fmt.Println("TEST MAIN")

	ret := m.Run()

	// setup

	os.Exit(ret)
}

func NewApiClient(token string) *client.McflyClient {
	return client.NewMcflyClient(server.URL, token)
}

func initDB() {
	mdb.Init()
}

func resetDB() {
	cleanupDB()
	seedDB()
}

func cleanupDB() {
	mdb.Clean()
}

func seedDB() {
	mdb.Seed()
}

func strPtr(v string) *string { return &v }
