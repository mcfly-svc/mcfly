package api

import (
	"net/http"
	"encoding/json"
)

func decodeRequestBodyJson(req *http.Request, v interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		return err
	}
	return nil
}