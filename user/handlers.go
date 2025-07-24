package user

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/omed0/go-hello-world/internal/api"
	"github.com/omed0/go-hello-world/internal/auth"
	"github.com/omed0/go-hello-world/internal/database"
	"github.com/omed0/go-hello-world/utils"
)

func HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
	}

	API := api.GetQueries()

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); utils.HandleRequestError(w, err, "Invalid JSON", http.StatusBadRequest) {
		return
	}

	isValid, validationErrors := validateUsername(params.Username)
	if !isValid {
		utils.RespondWithError(w, http.StatusBadRequest, validationErrors[0].Message)
		return
	}

	now := time.Now().UTC()
	user, err := API.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Username:  params.Username,
	})

	if utils.HandleRequestError(w, err, "Failed to create user", http.StatusInternalServerError) {
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := api.GetQueries().GetUserByAPIKey(r.Context(), apiKey)
	if utils.HandleRequestError(w, err, "Failed to get user", http.StatusInternalServerError) {
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

type failMessage struct {
	Message string `json:"message"`
}

// ValidateUsername checks if the username is valid according to the specified rules.
func validateUsername(username string) (bool, []failMessage) {

	if len(username) < 3 || len(username) > 25 {
		return false, []failMessage{{Message: "Username must be between 3 and 20 characters long"}}
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{1,28}[a-zA-Z0-9]$`)
	// This regex allows:
	// - Starts with a letter
	// - Contains letters, numbers, and underscores
	// - Ends with a letter or number
	// - Length between 3 and 30 characters

	if !validUsername.MatchString(username) {
		return false, []failMessage{{Message: "Username can only contain letters, numbers, and underscores, adn with correct format"}}
	}

	return true, nil
}
