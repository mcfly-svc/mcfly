package api

import (
	"net/http"
	"encoding/json"
)

type ApiResponse struct {
	Status 				string				`json:"status"`
}

type ApiError struct {
  Error         string        `json:"error"`
}

func writeErrorResponse(w http.ResponseWriter, errorMessage string) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusBadRequest)
  if err := json.NewEncoder(w).Encode(ApiError{errorMessage}); err != nil {
    panic(err)
  }
}

func writeSuccessResponse(w http.ResponseWriter) {  
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  if err := json.NewEncoder(w).Encode(ApiResponse{"success!"}); err != nil {
    panic(err)
  }
}
