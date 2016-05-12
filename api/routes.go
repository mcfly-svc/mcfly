package api

import (
	"net/http"
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
				RequestData:     LoginReq{},
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
				RequestData:       ProjectReq{},
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
				RequestData:       ProjectReq{},
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

		// projects
		/*Route{
			"ProjectPost",
			"POST",
			"/api/0/projects",
			handlers.ProjectPost,
		},
		Route{
			"ProjectsGet",
			"GET",
			"/api/0/projects",
			handlers.ProjectsGet,
		},
		Route{
			"ProjectGet",
			"GET",
			"/api/0/projects/{project_id}",
			handlers.ProjectGet,
		},
		Route{
			"ProjectDelete",
			"DELETE",
			"/api/0/projects/{project_id}",
			handlers.ProjectDelete,
		},
		*/
		// users
		/*Route{
		      "UserPost",
		      "POST",
		      "/api/0/users",
		      handlers.UserPost,
		  },
		  Route{
		      "UserGet",
		      "GET",
		      "/api/0/users/{user_id}",
		      handlers.UserGet,
		  },
		  Route{
		      "UserDelete",
		      "DELETE",
		      "/api/0/users/{user_id}",
		      handlers.UserDelete,
		  },*/
	}
}
