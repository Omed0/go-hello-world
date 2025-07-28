package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/omed0/go-hello-world/internal/auth"
	"github.com/omed0/go-hello-world/internal/database"
	"github.com/omed0/go-hello-world/models"
)

// HandlerCreateOrganization creates a new organization
func (api *ApiConfig) HandlerCreateOrganization(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	var params models.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Validate input
	params.Name = strings.TrimSpace(params.Name)
	if params.Name == "" {
		RespondWithError(w, http.StatusBadRequest, "Organization name is required")
		return
	}

	if len(params.Name) > 100 {
		RespondWithError(w, http.StatusBadRequest, "Organization name must be less than 100 characters")
		return
	}

	// Create organization
	createParams := database.CreateOrganizationParams{
		ID:   uuid.New(),
		Name: params.Name,
	}

	if params.Description != nil {
		createParams.Description.Valid = true
		createParams.Description.String = *params.Description
	}

	org, err := api.Queries.CreateOrganization(r.Context(), createParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create organization: "+err.Error())
		return
	}

	// Update user to be the owner of this organization
	_, err = api.Queries.UpdateUserRole(r.Context(), database.UpdateUserRoleParams{
		ID:   userID,
		Role: "owner",
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to set user as owner")
		return
	}

	// Set user's organization
	_, err = api.Queries.UpdateUserOrganization(r.Context(), database.UpdateUserOrganizationParams{
		ID:             userID,
		OrganizationID: uuid.NullUUID{UUID: org.ID, Valid: true},
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to assign user to organization")
		return
	}

	RespondWithJSON(w, http.StatusCreated, models.DatabaseOrganizationToOrganization(org))
}

// HandlerGetOrganization gets organization details
func (api *ApiConfig) HandlerGetOrganization(w http.ResponseWriter, r *http.Request) {
	orgIDStr := chi.URLParam(r, "orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid organization ID")
		return
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	// Get user to check permissions
	user, err := api.Queries.GetUserByID(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get user details")
		return
	}

	// Check if user has access to this organization
	if user.Role != "admin" && user.Role != "owner" {
		if !user.OrganizationID.Valid || user.OrganizationID.UUID != orgID {
			RespondWithError(w, http.StatusForbidden, "Access denied to this organization")
			return
		}
	}

	// Get organization
	org, err := api.Queries.GetOrganizationByID(r.Context(), orgID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Organization not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseOrganizationToOrganization(org))
}

// HandlerGetOrganizationUsers gets all users in an organization
func (api *ApiConfig) HandlerGetOrganizationUsers(w http.ResponseWriter, r *http.Request) {
	orgIDStr := chi.URLParam(r, "orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid organization ID")
		return
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	// Get user to check permissions
	user, err := api.Queries.GetUserByID(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get user details")
		return
	}

	// Check if user has access to this organization
	if user.Role != "admin" && user.Role != "owner" {
		if !user.OrganizationID.Valid || user.OrganizationID.UUID != orgID {
			RespondWithError(w, http.StatusForbidden, "Access denied to this organization")
			return
		}
	}

	// Get organization users
	users, err := api.Queries.GetUsersByOrganization(r.Context(), uuid.NullUUID{UUID: orgID, Valid: true})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get organization users")
		return
	}

	// Convert to response models
	var responseUsers []models.User
	for _, dbUser := range users {
		responseUsers = append(responseUsers, models.DatabaseUserRowToUser(dbUser))
	}

	RespondWithJSON(w, http.StatusOK, responseUsers)
}

// HandlerUpdateOrganization updates organization details
func (api *ApiConfig) HandlerUpdateOrganization(w http.ResponseWriter, r *http.Request) {
	orgIDStr := chi.URLParam(r, "orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid organization ID")
		return
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	// Get user to check permissions - only owners and admins can update organizations
	user, err := api.Queries.GetUserByID(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get user details")
		return
	}

	// Check permissions
	if user.Role != "admin" && user.Role != "owner" {
		RespondWithError(w, http.StatusForbidden, "Only owners and admins can update organizations")
		return
	}

	var params models.UpdateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	// Validate input
	if params.Name != nil {
		*params.Name = strings.TrimSpace(*params.Name)
		if *params.Name == "" {
			RespondWithError(w, http.StatusBadRequest, "Organization name cannot be empty")
			return
		}
		if len(*params.Name) > 100 {
			RespondWithError(w, http.StatusBadRequest, "Organization name must be less than 100 characters")
			return
		}
	}

	// Update organization
	updateParams := database.UpdateOrganizationParams{
		ID: orgID,
	}

	if params.Name != nil {
		updateParams.Column2 = *params.Name
	} else {
		updateParams.Column2 = "" // Empty string to preserve existing name
	}

	if params.Description != nil {
		updateParams.Description.Valid = true
		updateParams.Description.String = *params.Description
	}

	org, err := api.Queries.UpdateOrganization(r.Context(), updateParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to update organization")
		return
	}

	RespondWithJSON(w, http.StatusOK, models.DatabaseOrganizationToOrganization(org))
}

// HandlerDeleteOrganization soft deletes an organization
func (api *ApiConfig) HandlerDeleteOrganization(w http.ResponseWriter, r *http.Request) {
	orgIDStr := chi.URLParam(r, "orgId")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid organization ID")
		return
	}

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "User not found in context")
		return
	}

	// Get user to check permissions - only owners can delete organizations
	user, err := api.Queries.GetUserByID(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get user details")
		return
	}

	// Check permissions - only owners can delete organizations
	if user.Role != "owner" {
		RespondWithError(w, http.StatusForbidden, "Only organization owners can delete organizations")
		return
	}

	// Soft delete organization
	_, err = api.Queries.SoftDeleteOrganization(r.Context(), orgID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to delete organization")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
