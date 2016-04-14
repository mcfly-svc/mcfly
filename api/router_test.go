package api_test

import (
	"github.com/mikec/marsupi-api/api"
	"github.com/mikec/marsupi-api/apiutil"
  "github.com/mikec/marsupi-api/logging"

	"fmt"
  "io"
	"log"
	"os"
	"os/exec"
	"testing"
  "net/http/httptest"
)

var (
    server   			*httptest.Server
    reader   			io.Reader
    autil					apiutil.ApiUtil
)

/*type MockLogger struct {}

func (m MockLogger) Handler(h http.Handler, s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}*/

func init() {
    server = httptest.NewServer(api.NewRouter("postgres://localhost:5432/marsupi_test?sslmode=disable", logging.HttpRequestLogger{}))
		autil = apiutil.ApiUtil{server.URL}
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

