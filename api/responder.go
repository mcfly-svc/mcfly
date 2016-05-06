package api

import (
	"encoding/json"
	"net/http"

	"github.com/mikec/msplapi/logging"
)

type Responder struct {
	http.ResponseWriter
}

type ApiResponse struct {
	Status string `json:"status"`
}

func (r *Responder) Respond(v interface{}) {
	r.WriteCommonHeaders()
	r.WriteSuccessHeaders()
	r.WriteResponseData(v)
}

func (r *Responder) RespondWithError(apiErr *ApiError) {
	r.WriteCommonHeaders()
	r.WriteErrorHeaders()
	r.WriteResponseData(apiErr)
}

func (r *Responder) RespondWithServerError(err error) {
	logging.Panic(err)
	r.WriteCommonHeaders()
	r.WriteErrorHeaders()
	r.WriteResponseData(NewServerErr())
}

func (r *Responder) RespondWithSuccess() {
	r.Respond(ApiResponse{"success!"})
}

func (r *Responder) WriteCommonHeaders() {
	r.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func (r *Responder) WriteSuccessHeaders() {
	r.WriteHeader(http.StatusOK)
}

func (r *Responder) WriteErrorHeaders() {
	r.WriteHeader(http.StatusBadRequest)
}

func (r *Responder) WriteResponseData(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		logging.Panic(err)
	}
	r.Write(b)
}
