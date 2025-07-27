package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithJSON sends a JSON response with the specified status code and payload
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// RespondWithError sends an error response in JSON format
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		log.Printf("Server error: %s", msg)
	} else if code >= 400 {
		log.Printf("Client error: %s", msg)
	}

	RespondWithJSON(w, code, ErrorResponse{Error: msg})
}
