package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/omed0/go-hello-world/internal/database"
)

// User represents a user in the system
type User struct {
	ID               uuid.UUID  `json:"id"`
	Username         string     `json:"username"`
	Age              *int       `json:"age,omitempty"`
	Gender           *string    `json:"gender,omitempty"`
	Role             string     `json:"role"`
	OrganizationID   *uuid.UUID `json:"organization_id,omitempty"`
	OrganizationName *string    `json:"organization_name,omitempty"`
	APIKey           string     `json:"api_key"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Username       string     `json:"username" validate:"required,min=3,max=25"`
	Password       string     `json:"password" validate:"required,min=8"`
	Age            *int       `json:"age,omitempty" validate:"omitempty,min=13,max=120"`
	Gender         *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response for successful login
type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token,omitempty"` // For future JWT implementation
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=25"`
	Age      *int    `json:"age,omitempty" validate:"omitempty,min=13,max=120"`
	Gender   *string `json:"gender,omitempty" validate:"omitempty,oneof=male female other prefer_not_to_say"`
}

// DatabaseUserToUser converts a database user to a user model
func DatabaseUserToUser(dbUser database.User) User {
	user := User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Role:      dbUser.Role,
		APIKey:    dbUser.ApiKey,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	// Handle nullable fields
	if dbUser.Age.Valid {
		age := int(dbUser.Age.Int32)
		user.Age = &age
	}

	if dbUser.Gender.Valid {
		user.Gender = &dbUser.Gender.String
	}

	if dbUser.OrganizationID.Valid {
		user.OrganizationID = &dbUser.OrganizationID.UUID
	}

	return user
}

// DatabaseUserRowToUser converts a database user row to a user model
func DatabaseUserRowToUser(dbUser interface{}) User {
	switch v := dbUser.(type) {
	case database.GetUserByIDRow:
		return rowToUser(v.ID, v.Username, v.Role, v.ApiKey, v.CreatedAt, v.UpdatedAt, v.Age, v.Gender, v.OrganizationID, v.OrganizationName)
	case database.GetUserByUsernameRow:
		return rowToUser(v.ID, v.Username, v.Role, v.ApiKey, v.CreatedAt, v.UpdatedAt, v.Age, v.Gender, v.OrganizationID, v.OrganizationName)
	case database.GetUserByAPIKeyRow:
		return rowToUser(v.ID, v.Username, v.Role, v.ApiKey, v.CreatedAt, v.UpdatedAt, v.Age, v.Gender, v.OrganizationID, v.OrganizationName)
	case database.GetUserByUsernameAndPasswordRow:
		return rowToUser(v.ID, v.Username, v.Role, v.ApiKey, v.CreatedAt, v.UpdatedAt, v.Age, v.Gender, v.OrganizationID, v.OrganizationName)
	default:
		// Fallback to empty user if unknown type
		return User{}
	}
}

// Helper function to convert row data to User
func rowToUser(id uuid.UUID, username string, role string, apiKey string, createdAt time.Time, updatedAt time.Time, age sql.NullInt32, gender sql.NullString, organizationID uuid.NullUUID, organizationName sql.NullString) User {
	user := User{
		ID:        id,
		Username:  username,
		Role:      role,
		APIKey:    apiKey,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	// Handle nullable fields
	if age.Valid {
		ageInt := int(age.Int32)
		user.Age = &ageInt
	}

	if gender.Valid {
		user.Gender = &gender.String
	}

	if organizationID.Valid {
		user.OrganizationID = &organizationID.UUID
	}

	if organizationName.Valid {
		user.OrganizationName = &organizationName.String
	}

	return user
}
