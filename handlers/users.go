package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/omed0/go-hello-world/internal/auth"
	"github.com/omed0/go-hello-world/internal/database"
	"github.com/omed0/go-hello-world/models"
)

func (api *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); HandleRequestError(w, err, "Invalid JSON", http.StatusBadRequest) {
		return
	}

	if valid, errMsg := validateUsername(params.Username); !valid {
		RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	user, err := api.Queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:       uuid.New(),
		Username: params.Username,
	})
	if HandleRequestError(w, err, "Failed to create user", http.StatusInternalServerError) {
		return
	}

	RespondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
}

func (api ApiConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := api.Queries.GetUserByAPIKey(r.Context(), apiKey)
	if HandleRequestError(w, err, "Failed to get user", http.StatusInternalServerError) {
		return
	}
	RespondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}

func validateUsername(username string) (bool, string) {
	if len(username) < 3 || len(username) > 25 {
		return false, "Username must be between 3 and 25 characters long"
	}
	validUsername := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{1,23}[a-zA-Z0-9]$`)
	if !validUsername.MatchString(username) {
		return false, "Username can only contain letters, numbers, and underscores, and must start with a letter and end with a letter or number"
	}
	return true, ""
}
