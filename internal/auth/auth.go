package auth

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// Context key type for user ID - use unexported type to prevent collisions
type contextKey string

const userIDKey contextKey = "userID"

// Compile regex once at package level for better performance
var apiKeyRegex = regexp.MustCompile(`^APIKEY\s+([a-zA-Z0-9]+)$`)

// Common errors - Define clear, user-friendly error messages
var (
	ErrMissingAuthHeader   = errors.New("missing authorization header")
	ErrInvalidAPIKeyFormat = errors.New("invalid API key format, expected: APIKEY <key>")
	ErrInvalidAPIKey       = errors.New("invalid API key")
	ErrUserNotFound        = errors.New("user not found")
	ErrUnauthorized        = errors.New("unauthorized access, please provide a valid API key")
	ErrInternalServer      = errors.New("internal server error, please try again later")
)

// AuthenticateFunc is a function type that validates an API key and returns a user ID
type AuthenticateFunc func(apiKey string) (uuid.UUID, error)

// GetAPIKey extracts and validates the API key from the Authorization header
func GetAPIKey(h http.Header) (string, error) {
	apiKey := strings.TrimSpace(h.Get("Authorization"))
	if apiKey == "" {
		return "", ErrMissingAuthHeader
	}

	matches := apiKeyRegex.FindStringSubmatch(apiKey)
	if len(matches) < 2 {
		return "", ErrInvalidAPIKeyFormat
	}

	return matches[1], nil
}

// SetUserIDInContext sets user ID in request context
func SetUserIDInContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext retrieves user ID from request context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New(ErrUnauthorized.Error())
	}
	return userID, nil
}

// Middleware creates an authentication middleware
func Middleware(authenticateFunc AuthenticateFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract API key from header
			apiKey, err := GetAPIKey(r.Header)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Validate API key and get user ID
			userID, err := authenticateFunc(apiKey)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Set user ID in context
			ctx := SetUserIDInContext(r.Context(), userID)
			r = r.WithContext(ctx)

			// Call next handler
			next.ServeHTTP(w, r)
		})
	}
}
