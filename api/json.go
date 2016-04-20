package api

import (
	"net/http"
	"encoding/json"
)

func DecodeRequest(req *http.Request, v interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func DecodeResponse(res *http.Response, v interface{}) error {
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}
	return nil
}