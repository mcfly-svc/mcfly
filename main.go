package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mikec/marsupi-api/api"
	"github.com/mikec/marsupi-api/logging"
	"github.com/mikec/marsupi-api/provider"
)

func main() {

	args := os.Args[1:]
	if len(args) > 0 {
		RunCommands(args)
		os.Exit(1)
		return
	}

	github := provider.GitHub{}
	dropbox := provider.Dropbox{}

	authProviders := make(map[string]provider.AuthProvider)
	authProviders[github.Key()] = &github
	authProviders[dropbox.Key()] = &dropbox

	router := api.NewRouter(
		"postgres://localhost:5432/marsupi_test?sslmode=disable",
		logging.HttpRequestLogger{},
		authProviders,
	)

	log.Fatal(http.ListenAndServe(":8081", router))

}
