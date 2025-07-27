package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/omed0/go-hello-world/internal/auth"
	"github.com/omed0/go-hello-world/internal/database"
	"github.com/omed0/go-hello-world/models"
)

// Compile regex once at package level for better performance
var validUsernameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{1,23}[a-zA-Z0-9]$`)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username string `json:"username"`
}

func (api *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var params CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Trim whitespace and validate
	params.Username = strings.TrimSpace(params.Username)
	if valid, errMsg := validateUsername(params.Username); !valid {
		RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	user, err := api.Queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:       uuid.New(),
		Username: params.Username,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
}

func (api *ApiConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := api.Queries.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}

// validateUsername validates the username according to business rules
func validateUsername(username string) (bool, string) {
	if len(username) < 3 || len(username) > 25 {
		return false, "Username must be between 3 and 25 characters long"
	}

	if !validUsernameRegex.MatchString(username) {
		return false, "Username can only contain letters, numbers, and underscores, and must start with a letter and end with a letter or number"
	}

	return true, ""
}
