package api

import (
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/mikec/msplapi/models"
	"github.com/mikec/msplapi/mq"
	"github.com/mikec/msplapi/provider"
)

type Handlers struct {
	DB              models.Datastore
	MessageChannel  mq.MessageChannel
	GenerateToken   func() string
	AuthProviders   map[string]provider.AuthProvider
	SourceProviders map[string]provider.SourceProvider
}

type HandlerOptions struct {
	AuthRequired      bool
	RequestData       interface{}
	ProjectContext    bool
	UseSourceProvider bool
	UseAuthProvider   bool
	After             func(*Responder, *RequestContext)
}

type RequestContext struct {
	CurrentUser         *models.User
	RequestData         interface{}
	SourceProvider      *provider.SourceProvider
	SourceProviderToken *models.ProviderAccessToken
	Project             *models.Project
	AuthProvider        *provider.AuthProvider
}

type SourceProviderRequest interface {
	SourceProvider() string
}

type AuthProviderRequest interface {
	AuthProvider() string
}

type ProjectContextRequest interface {
	ProjectHandle() string
}

// MakeHandlerFunc returns a handler function based on options provided. It executes the
// HandlerOptions.After function if all options run successfully.
func (handlers *Handlers) MakeHandlerFunc(opts HandlerOptions) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := &RequestContext{}

		vars := mux.Vars(req)

		r := &Responder{w, req}

		if opts.AuthRequired {
			ctx.CurrentUser = r.ValidateAuthorization(handlers.DB)
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

		if opts.ProjectContext {
			opts.UseSourceProvider = true
		}

		var sp string
		if opts.UseSourceProvider {
			if vars["provider"] != "" {
				sp = vars["provider"]
			} else {
				sp = ctx.RequestData.(SourceProviderRequest).SourceProvider()
			}
			sourceProvider := handlers.SourceProviders[sp]
			if sourceProvider == nil {
				// TODO: change these errors to provider type specific
				r.RespondWithError(NewUnsupportedProviderErr(sp))
				return
			}
			ctx.SourceProvider = &sourceProvider
			if opts.AuthRequired {
				spToken, err := handlers.DB.GetProviderTokenForUser(ctx.CurrentUser, sp)
				if err != nil {
					if err == models.ErrNotFound {
						r.RespondWithError(NewProviderTokenNotFoundErr(sp))
						return
					}
					r.RespondWithServerError(err)
					return
				}
				ctx.SourceProviderToken = spToken
			}
		}

		if opts.ProjectContext {
			projectHandle := ctx.RequestData.(ProjectContextRequest).ProjectHandle()
			project, err := handlers.DB.GetProject(projectHandle, sp)
			if err != nil {
				if err == models.ErrNotFound {
					r.RespondWithError(NewNotFoundErr("project", projectHandle))
					return
				}
				r.RespondWithServerError(err)
				return
			}
			ctx.Project = project
		}

		if opts.UseAuthProvider {
			ap := ctx.RequestData.(AuthProviderRequest).AuthProvider()
			authProvider := handlers.AuthProviders[ap]
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
