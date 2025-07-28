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
var usernameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{1,23}[a-zA-Z0-9]$`)

// HandlerLogin handles user login with username and password
func (api *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var params models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Trim whitespace
	params.Username = strings.TrimSpace(params.Username)
	params.Password = strings.TrimSpace(params.Password)

	// Validate input
	if params.Username == "" || params.Password == "" {
		RespondWithError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	// Hash the password for comparison
	passwordHash, err := auth.HashPassword(params.Password, nil)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Authentication failed")
		return
	}

	// Get user by username and password
	user, err := api.Queries.GetUserByUsernameAndPassword(r.Context(), database.GetUserByUsernameAndPasswordParams{
		Username:     params.Username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// Return user details and API key
	response := models.LoginResponse{
		User: models.DatabaseUserRowToUser(user),
	}

	RespondWithJSON(w, http.StatusOK, response)
}

// HandlerCreateUser creates a new user with enhanced fields
func (api *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var params models.CreateUserRequest

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

	// Validate password strength
	if err := auth.ValidatePasswordStrength(params.Password); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Hash password
	passwordHash, err := auth.HashPassword(params.Password, nil)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Validate age if provided
	if params.Age != nil && (*params.Age < 13 || *params.Age > 120) {
		RespondWithError(w, http.StatusBadRequest, "Age must be between 13 and 120")
		return
	}

	// Validate gender if provided
	if params.Gender != nil {
		validGenders := []string{"male", "female", "other", "prefer_not_to_say"}
		valid := false
		for _, g := range validGenders {
			if *params.Gender == g {
				valid = true
				break
			}
		}
		if !valid {
			RespondWithError(w, http.StatusBadRequest, "Invalid gender value")
			return
		}
	}

	// Create user in database
	createParams := database.CreateUserWithPasswordParams{
		ID:           uuid.New(),
		Username:     params.Username,
		PasswordHash: passwordHash,
	}

	// Handle optional fields
	if params.Age != nil {
		createParams.Age.Valid = true
		createParams.Age.Int32 = int32(*params.Age)
	}

	if params.Gender != nil {
		createParams.Gender.Valid = true
		createParams.Gender.String = *params.Gender
	}

	// Set default role as 'user'
	createParams.Column6 = "user"

	// Handle organization ID
	if params.OrganizationID != nil {
		createParams.Column7 = *params.OrganizationID
	}

	user, err := api.Queries.CreateUserWithPassword(r.Context(), createParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
}

func (api *ApiConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	// Get user by ID instead of API key since we already have it from middleware
	user, err := api.Queries.GetUserByID(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get user: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseUserRowToUser(user))
}

// HandlerUpdateUser updates user information
func (api *ApiConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	var params models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Validate username if provided
	if params.Username != nil {
		*params.Username = strings.TrimSpace(*params.Username)
		if valid, errMsg := validateUsername(*params.Username); !valid {
			RespondWithError(w, http.StatusBadRequest, errMsg)
			return
		}
	}

	// Validate age if provided
	if params.Age != nil && (*params.Age < 13 || *params.Age > 120) {
		RespondWithError(w, http.StatusBadRequest, "Age must be between 13 and 120")
		return
	}

	// Validate gender if provided
	if params.Gender != nil {
		validGenders := []string{"male", "female", "other", "prefer_not_to_say"}
		valid := false
		for _, g := range validGenders {
			if *params.Gender == g {
				valid = true
				break
			}
		}
		if !valid {
			RespondWithError(w, http.StatusBadRequest, "Invalid gender value")
			return
		}
	}

	// Create update parameters
	updateParams := database.UpdateUserParams{
		ID: userID,
	}

	if params.Username != nil {
		updateParams.Username = *params.Username
	}

	if params.Age != nil {
		updateParams.Age.Valid = true
		updateParams.Age.Int32 = int32(*params.Age)
	}

	if params.Gender != nil {
		updateParams.Gender.Valid = true
		updateParams.Gender.String = *params.Gender
	}

	// Update user
	user, err := api.Queries.UpdateUser(r.Context(), updateParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}
