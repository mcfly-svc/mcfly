package api_test

import (
	"github.com/mikec/marsupi-api/api"
	"github.com/mikec/marsupi-api/client"
  "github.com/mikec/marsupi-api/logging"

	"fmt"
  "io"
	"log"
	"os"
	"os/exec"
	"testing"
	"net/http"
  "net/http/httptest"
)

var (
    server   				*httptest.Server
    reader   				io.Reader
    apiClient				client.Client
)

/*type MockLogger struct {}

func (m MockLogger) Handler(h http.Handler, s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}*/

func init() {
    server = httptest.NewServer(api.NewRouter("postgres://localhost:5432/marsupi_test?sslmode=disable", logging.HttpRequestLogger{}))
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
	runHelperScript("recreate.sh")
}

func cleanupDB() {
	runHelperScript("clean.sh")
}

func runHelperScript(sh string) {
	out, err := exec.Command("/bin/sh", fmt.Sprintf("../db/helpers/%s", sh)).Output()
	if err != nil {
		log.Fatal(fmt.Sprintf("%s failed: ", sh), err)
	}
	fmt.Printf("%s", out)
}

type EndpointTester struct {
	Test 				*testing.T
	Endpoint 		client.EntityEndpoint
}

func (self *EndpointTester) Create(JSON string) *http.Response {
	res, err := self.Endpoint.Create(JSON)
  if err != nil {
    self.Test.Error(err)
  }
  return res
}

func (self *EndpointTester) GetAll() []interface{} {
	res, err := self.Endpoint.GetAll()
  if err != nil {
    self.Test.Error(err)
  }
  entities, err := self.Endpoint.Decoder.DecodeArrayResponse(res)
  if err != nil {
  	self.Test.Error(err)
  }
  return entities
}

func (self *EndpointTester) GetRes(ID int64) *http.Response {
	res, err := self.Endpoint.Get(ID)
  if err != nil {
    self.Test.Error(err)
  }
  return res
}

func (self *EndpointTester) Get(ID int64) interface{} {
	res := self.GetRes(ID)
  entity, err := self.Endpoint.Decoder.DecodeResponse(res)
  if err != nil {
  	self.Test.Error(err)
  }
  return entity
}

func (self *EndpointTester) Delete(ID int64) *http.Response {
	res, err := self.Endpoint.Delete(ID)
	if err != nil {
		self.Test.Error(err)
	}
	return res
}

