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
        Route{
            "ProjectsPost",
            "POST",
            "/api/0/projects",
            handlers.ProjectsPost,
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
            "ProjectsDelete",
            "DELETE",
            "/api/0/projects/{project_id}",
            handlers.ProjectsDelete,
        },
    }
}
