package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/omed0/go-hello-world/internal/database"
)

// Organization represents an organization in the system
type Organization struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// CreateOrganizationRequest represents the request body for creating an organization
type CreateOrganizationRequest struct {
	Name        string                 `json:"name" validate:"required,min=2,max=100"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=500"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
}

// UpdateOrganizationRequest represents the request body for updating an organization
type UpdateOrganizationRequest struct {
	Name        *string                `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string                `json:"description,omitempty" validate:"omitempty,max=500"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
}

// DatabaseOrganizationToOrganization converts a database organization to an organization model
func DatabaseOrganizationToOrganization(dbOrg database.Organization) Organization {
	org := Organization{
		ID:        dbOrg.ID,
		Name:      dbOrg.Name,
		CreatedAt: dbOrg.CreatedAt,
		UpdatedAt: dbOrg.UpdatedAt,
	}

	// Handle nullable description
	if dbOrg.Description.Valid {
		org.Description = &dbOrg.Description.String
	}

	// Handle JSONB settings
	if dbOrg.Settings.Valid {
		// Note: In a real implementation, you'd unmarshal the JSON
		// For now, we'll leave it empty
		org.Settings = make(map[string]interface{})
	}

	return org
}

// DatabaseOrganizationsToOrganizations converts a slice of database organizations to organization models
func DatabaseOrganizationsToOrganizations(dbOrgs []database.Organization) []Organization {
	orgs := make([]Organization, len(dbOrgs))
	for i, dbOrg := range dbOrgs {
		orgs[i] = DatabaseOrganizationToOrganization(dbOrg)
	}
	return orgs
}
