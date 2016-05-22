package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikec/msplapi/config"
	"github.com/mikec/msplapi/mq"

	"github.com/mikec/msplapi/models"
	"github.com/mikec/msplapi/provider"
)

type RequestLogger interface {
	Handler(http.Handler, string) http.Handler
}

func NewRouter(
	cfg *config.Config,
	logger RequestLogger,
	msgChannel mq.MessageChannel,
	generateToken func() string,
	authProviders map[string]provider.AuthProvider,
	sourceProviders map[string]provider.SourceProvider,
) *mux.Router {

	db, err := models.NewDB(cfg.DatabaseUrl, cfg.DatabaseName, cfg.DatabaseUseSSL)
	if err != nil {
		log.Panic(err)
	}

	handlers := &Handlers{
		DB:              db,
		MessageChannel:  msgChannel,
		GenerateToken:   generateToken,
		AuthProviders:   authProviders,
		SourceProviders: sourceProviders,
	}
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
