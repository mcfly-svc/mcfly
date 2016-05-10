package api

import (
	"net/http"
	"reflect"

	"github.com/mikec/msplapi/models"
	"github.com/mikec/msplapi/provider"
)

type Handlers struct {
	db              models.Datastore
	generateToken   func() string
	authProviders   map[string]provider.AuthProvider
	sourceProviders map[string]provider.SourceProvider
}

type HandlerOptions struct {
	AuthRequired      bool
	RequestData       interface{}
	UseSourceProvider bool
	UseAuthProvider   bool
	After             func(*Responder, *RequestContext)
}

type RequestContext struct {
	CurrentUser    *models.User
	RequestData    interface{}
	SourceProvider *provider.SourceProvider
	AuthProvider   *provider.AuthProvider
}

type SourceProviderRequest interface {
	SourceProvider() string
}

type AuthProviderRequest interface {
	AuthProvider() string
}

// MakeHandlerFunc returns a handler function based on options provided. It executes the
// HandlerOptions.After function if all options run successfully.
func (handlers *Handlers) MakeHandlerFunc(opts HandlerOptions) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := &RequestContext{}

		r := &Responder{w, req}

		if opts.AuthRequired {
			ctx.CurrentUser = r.ValidateAuthorization(handlers.db)
			if ctx.CurrentUser == nil {
				return
			}
		}

		if opts.RequestData != nil {
			reqDataType := reflect.TypeOf(opts.RequestData)
			ctx.RequestData = reflect.New(reqDataType).Interface()

			decodeErr := r.DecodeRequest(ctx.RequestData)
			if decodeErr != nil {
				return
			}
			reqValid := r.ValidateRequestData(ctx.RequestData)
			if !reqValid {
				return
			}
		}

		if opts.UseSourceProvider {
			sp := ctx.RequestData.(SourceProviderRequest).SourceProvider()
			sourceProvider := handlers.sourceProviders[sp]
			if sourceProvider == nil {
				// TODO: change these errors to provider type specific
				r.RespondWithError(NewUnsupportedProviderErr(sp))
				return
			}
			ctx.SourceProvider = &sourceProvider
		}

		if opts.UseAuthProvider {
			ap := ctx.RequestData.(AuthProviderRequest).AuthProvider()
			authProvider := handlers.authProviders[ap]
			if authProvider == nil {
				// TODO: change these errors to provider type specific
				r.RespondWithError(NewUnsupportedProviderErr(ap))
				return
			}
			ctx.AuthProvider = &authProvider
		}

		opts.After(r, ctx)
	}
}
