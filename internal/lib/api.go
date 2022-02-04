package lib

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is an api error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse is an empty api response
type SuccessResponse struct {
	Success bool `json:"success"`
}

// WriteResponseError writes an error with json body
func WriteErrorResponse(w http.ResponseWriter, status int, error string) {
	w.Header().Set("Content-Type", "application/json")
	e := &ErrorResponse{Error: error}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(e)
}
