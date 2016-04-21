package api

import (
  "fmt"
  "log"
	"net/http"
	"encoding/json"

  //"github.com/mikec/marsupi-api/models"
)

type Responder struct {
  http.ResponseWriter
}

type ApiResponse struct {
	Status 				string				`json:"status"`
}

type ApiError struct {
  Error         string        `json:"error"`
}

func (r *Responder) Respond(v interface{}) {
  r.WriteCommonHeaders()
  r.WriteSuccessHeaders()
  r.WriteResponseData(v)
}

func (r *Responder) RespondWithError(v interface{}) {
  r.WriteCommonHeaders()
  r.WriteErrorHeaders()
  r.WriteErrorData(v)
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
    log.Fatal(err)
  }
  r.Write(b)
}

func (r *Responder) WriteErrorData(v interface{}) {
  switch err := v.(type) {
  case string:
    v = ApiError{err}
  case ApiError:
    v = err
  default:
    log.Fatal(fmt.Errorf("Unexpected type %T", v))
  }

  b, err := json.Marshal(v)
  if err != nil {
    log.Fatal(err)
  }
  r.Write(b)
}
