package handlers

import (
	"net/http"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// HandlerReadiness handles health check requests
// This endpoint can be used by load balancers and monitoring systems
func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	if !IsDBConnected() {
		RespondWithJSON(w, http.StatusServiceUnavailable, HealthResponse{
			Status:  "unhealthy",
			Message: "Database connection is not available",
		})
		return
	}

	RespondWithJSON(w, http.StatusOK, HealthResponse{
		Status:  "healthy",
		Message: "Service is running properly",
	})
}

// HandlerErr handles error testing requests
// This endpoint is useful for testing error handling and logging
func HandlerErr(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusInternalServerError, "This is a test error endpoint")
}
