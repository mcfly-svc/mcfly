package main

import (
    "log"
    "net/http"

    "github.com/mikec/marsupi-api/api"
)

func main() {

  router := api.NewRouter("postgres://localhost:5432/marsupi_test?sslmode=disable")

  log.Fatal(http.ListenAndServe(":8080", router))

}
