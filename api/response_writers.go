package api

import (
  "fmt"
  "log"
	"net/http"
	"encoding/json"
)

type ApiResponse struct {
	Status 				string				`json:"status"`
}

type ApiError struct {
  Error         string        `json:"error"`
}

func writeErrorResponse(w http.ResponseWriter, apiError interface{}) {

  switch v := apiError.(type) {
  case string:
    apiError = ApiError{v}
  case ApiError:
    apiError = v
  default:
    log.Fatal(fmt.Errorf("Unexpected type %T", apiError))
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusBadRequest)
  if err := json.NewEncoder(w).Encode(apiError); err != nil {
    log.Fatal(err)
  }
}

func writeSuccessResponse(w http.ResponseWriter) {  
  writeResponse(w, ApiResponse{"success!"})
}

func writeResponse(w http.ResponseWriter, v interface{}) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(w).Encode(v); err != nil {
    log.Fatal(err)
  }
}
