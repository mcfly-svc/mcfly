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

        // projects
        Route{
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

        // users
        Route{
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
        },
    }
}
