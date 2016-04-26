package api

import (
    "log"
    "net/http"
    "math"

    "github.com/gorilla/mux"
    "github.com/RangelReale/osin"

    "github.com/mikec/marsupi-api/models"
)

type Handlers struct {
    db models.Datastore
    osinServer osin.Server
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

    store := &OAuthStorage{db}

    config := osin.NewServerConfig()
    config.AllowGetAccessRequest = true
    config.AllowedAccessTypes = osin.AllowedAccessType{
        osin.PASSWORD,
    }
    config.ErrorStatusCode = http.StatusBadRequest
    config.AccessExpiration = math.MaxInt32 // never expire

    osinServer := *osin.NewServer(config, store)

    handlers := &Handlers{db, osinServer, github}
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
