package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

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

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		log.Printf("Server error: %s", msg)
	}
	RespondWithJSON(w, code, ErrorResponse{Error: msg})
}
