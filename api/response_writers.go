package api

import (
	"net/http"
	"encoding/json"
)

type ApiResponse struct {
	Status 				string				`json:"status"`
}

func writeServerError(w http.ResponseWriter) {
	writeErrorMessage(w, "Server error")
}

func writeErrorMessage(w http.ResponseWriter, errorMessage string) {
	http.Error(w, errorMessage, http.StatusBadRequest)
}

func writeSuccessResponse(w http.ResponseWriter) {  
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  successResponse := ApiResponse{
  	Status: "success!",
  }
  if err := json.NewEncoder(w).Encode(successResponse); err != nil {
    panic(err)
  }
}
