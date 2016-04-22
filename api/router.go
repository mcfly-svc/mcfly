package api

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"

    "github.com/mikec/marsupi-api/models"
)

type Handlers struct {
    db models.Datastore
    github GitHubClient
}

type RequestLogger interface {
    Handler(http.Handler, string) http.Handler
}

type GitHubClient interface {
    GetAuthenticatedUser(string) (*models.User, error)
}

func NewRouter(dbUrl string, logger RequestLogger, github GitHubClient) *mux.Router {

    db, err := models.NewDB(dbUrl)
    if err != nil {
        log.Panic(err)
    }
    handlers := &Handlers{db, github}
    log.Printf("Connected to postgres")

    router := mux.NewRouter().StrictSlash(true)

    routes := AllRoutes(handlers)

    for _, route := range routes {

        var handler http.Handler

        handler = logger.Handler(route.HandlerFunc, route.Name)

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }

    return router
}
