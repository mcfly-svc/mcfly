package main

import (
    "log"
    "net/http"

    "github.com/mikec/marsupi-api/api"
    "github.com/mikec/marsupi-api/logging"
)

func main() {

  router := api.NewRouter("postgres://localhost:5432/marsupi_test?sslmode=disable", logging.HttpRequestLogger{})

  log.Fatal(http.ListenAndServe(":8080", router))

}
