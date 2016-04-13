package api_test

import (
	"github.com/mikec/marsupi-api/api"

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
    projectsUrl		string
)

func init() {
    server = httptest.NewServer(api.NewRouter("postgres://localhost:5432/marsupi_test?sslmode=disable")) 

    projectsUrl = fmt.Sprintf("%s/api/0/projects", server.URL)
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

