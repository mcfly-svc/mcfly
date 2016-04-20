package apiutil

import (
	"net/http"

	"github.com/mikec/marsupi-api/models"
	"github.com/mikec/marsupi-api/api"
)

type EntityDecoder interface {
	DecodeRequest(*http.Request) (interface{}, error)
	DecodeResponse(res *http.Response) (interface{}, error)
	DecodeArrayRequest(req *http.Request) ([]interface{}, error)
	DecodeArrayResponse(res *http.Response) ([]interface{}, error)
}

type ProjectDecoder struct {}

func (self ProjectDecoder) DecodeRequest(req *http.Request) (interface{}, error) {
    var p models.Project
    if err := api.DecodeRequest(req, &p); err != nil {
        return nil, err
    }
    return p, nil
}

func (self ProjectDecoder) DecodeResponse(res *http.Response) (interface{}, error) {
    var p models.Project
    if err := api.DecodeResponse(res, &p); err != nil {
        return nil, err
    }
    return p, nil
}

func (self ProjectDecoder) DecodeArrayRequest(req *http.Request) ([]interface{}, error) {
    var projects []models.Project
    if err := api.DecodeRequest(req, &projects); err != nil {
        return nil, err
    }
    ret := make([]interface{}, len(projects))
    for i, d := range projects {
        ret[i] = d
    }
    return ret, nil
}

func (self ProjectDecoder) DecodeArrayResponse(res *http.Response) ([]interface{}, error) {
    var projects []models.Project
    if err := api.DecodeResponse(res, &projects); err != nil {
        return nil, err
    }
    ret := make([]interface{}, len(projects))
    for i, d := range projects {
        ret[i] = d
    }
    return ret, nil
}
