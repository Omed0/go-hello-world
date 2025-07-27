package auth

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

// Compile regex once at package level for better performance
var apiKeyRegex = regexp.MustCompile(`^APIKEY\s+([a-zA-Z0-9]+)$`)

// Common errors
var (
	ErrMissingAuthHeader   = errors.New("missing authorization header")
	ErrInvalidAPIKeyFormat = errors.New("invalid API key format, expected: APIKEY <key>")
)

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
