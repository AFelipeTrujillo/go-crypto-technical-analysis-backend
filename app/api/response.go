package api

import (
	"encoding/json"
	"net/http"
)

type APIErrorResponse struct {
	Error string `json:"error"`
}

func OKResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Status Internal Server Error")
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIErrorResponse{Error: message})
}
