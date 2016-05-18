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
	"github.com/mikec/msplapi/provider/mockprovider"
	"github.com/streadway/amqp"

	"io"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	server    *httptest.Server
	reader    io.Reader
	apiClient *client.MsplClient
	cfg       *config.Config
	database  *models.DB
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
	cfg = GetTestConfig()
	dbs, err := models.NewDB(cfg.DatabaseUrl)

	database = dbs
	if err != nil {
		log.Fatal(err)
	}
	createDB()

	jabroni := mockprovider.MockProvider{}

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
	apiClient = client.NewMsplClient(server.URL, "")
}

func TestMain(m *testing.M) {

	ret := m.Run()

	// setup

	os.Exit(ret)
}

func NewApiClient(token string) *client.MsplClient {
	return client.NewMsplClient(server.URL, token)
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

func strPtr(v string) *string { return &v }
