package api

import (
	"net/http"

	"github.com/mikec/msplapi/api/apidata"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func AllRoutes(handlers *Handlers) Routes {
	return Routes{

		// auth
		Route{
			"Auth",
			"POST",
			"/api/0/login",
			handlers.MakeHandlerFunc(HandlerOptions{
				RequestData:     apidata.LoginReq{},
				UseAuthProvider: true,
				After:           handlers.Login,
			}),
		},

		// projects
		Route{
			"PostProject",
			"POST",
			"/api/0/projects",
			handlers.MakeHandlerFunc(HandlerOptions{
				AuthRequired:      true,
				RequestData:       apidata.ProjectReq{},
				UseSourceProvider: true,
				After:             handlers.PostProject,
			}),
		},

		Route{
			"GetProviderProjects",
			"GET",
			"/api/0/{provider}/projects",
			handlers.MakeHandlerFunc(HandlerOptions{
				AuthRequired:      true,
				UseSourceProvider: true,
				After:             handlers.GetProviderProjects,
			}),
		},

		Route{
			"GetProjects",
			"GET",
			"/api/0/projects",
			handlers.MakeHandlerFunc(HandlerOptions{
				AuthRequired: true,
				After:        handlers.GetProjects,
			}),
		},

		Route{
			"DeleteProject",
			"DELETE",
			"/api/0/projects",
			handlers.MakeHandlerFunc(HandlerOptions{
				AuthRequired:      true,
				RequestData:       apidata.ProjectReq{},
				UseSourceProvider: true,
				After:             handlers.DeleteProject,
			}),
		},

		Route{
			"ProjectUpdateWebhook",
			"POST",
			"/api/0/webhooks/{provider}/project-update",
			handlers.MakeHandlerFunc(HandlerOptions{
				UseSourceProvider: true,
				After:             handlers.ProjectUpdateWebhook,
			}),
		},

		Route{
			"SaveBuild",
			"POST",
			"/api/0/builds",
			handlers.MakeHandlerFunc(HandlerOptions{
				// TODO: will be called by another service, and that call
				// should probably be authenticated? the service won't have
				// the user's access token
				RequestData:       apidata.BuildReq{},
				UseSourceProvider: true,
				After:             handlers.SaveBuild,
			}),
		},

		Route{
			"Deploy",
			"POST",
			"/api/0/deploy",
			handlers.MakeHandlerFunc(HandlerOptions{
				AuthRequired:   true,
				RequestData:    apidata.DeployReq{},
				ProjectContext: true,
				After:          handlers.Deploy,
			}),
		},
	}
}
