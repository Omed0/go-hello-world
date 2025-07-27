package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// HandlerReadiness handles health check requests
func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

// HandlerErr handles error testing requests
func HandlerErr(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusBadRequest, "Something went wrong")
}

// HandleRequestError is a helper function to handle errors consistently
// Deprecated: Use direct error handling in handlers for better clarity
func HandleRequestError(w http.ResponseWriter, err error, msg string, code int) bool {
	if err != nil {
		log.Printf("Request error: %s - %v", msg, err)
		RespondWithError(w, code, fmt.Sprintf("%s: %v", msg, err))
		return true
	}
	return false
}
