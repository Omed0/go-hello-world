package auth

import (
	"errors"
	"net/http"
	"regexp"
)

var apiKeyRegex = regexp.MustCompile(`^APIKEY\s+([a-zA-Z0-9]+)$`)

func GetAPIKey(h http.Header) (string, error) {
	apiKey := h.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("missing authorization header")
	}

	matches := apiKeyRegex.FindStringSubmatch(apiKey)
	if len(matches) < 2 {
		return "", errors.New("invalid API key format")
	}

	return matches[1], nil
}
