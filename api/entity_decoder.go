package api

import "net/http"

type EntityDecoder interface {
	DecodeRequest(*http.Request) (interface{}, error)
	DecodeResponse(res *http.Response) (interface{}, error)
	DecodeArrayRequest(req *http.Request) ([]interface{}, error)
	DecodeArrayResponse(res *http.Response) ([]interface{}, error)
}
