package handlers

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

// validateUsername validates username format and requirements
func validateUsername(username string) (bool, string) {
	username = strings.TrimSpace(username)

	if len(username) < 3 {
		return false, "Username must be at least 3 characters long"
	}

	if len(username) > 50 {
		return false, "Username must be less than 50 characters"
	}

	// Allow alphanumeric characters, underscores, and hyphens
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
	if !matched {
		return false, "Username can only contain letters, numbers, underscores, and hyphens"
	}

	return true, ""
}

// parseLimit parses and validates limit parameter
func parseLimit(limitStr string) int {
	if limitStr == "" {
		return defaultLimit
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return defaultLimit
	}

	if limit > maxLimit {
		return maxLimit
	}

	return limit
}

func (e *ValidationError) Error() string {
	return e.Message
}
