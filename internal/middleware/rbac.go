package middleware

import (
	"net/http"
	"slices"

	"github.com/omed0/go-hello-world/handlers"
	"github.com/omed0/go-hello-world/internal/auth"
)

// Role constants
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
	RoleMod   = "moderator"
	RoleOwner = "owner"
)

// Permission levels
const (
	PermissionRead   = "read"
	PermissionWrite  = "write"
	PermissionDelete = "delete"
	PermissionAdmin  = "admin"
)

// RolePermissions defines what each role can do
var RolePermissions = map[string][]string{
	RoleUser:  {PermissionRead, PermissionWrite},
	RoleMod:   {PermissionRead, PermissionWrite, PermissionDelete},
	RoleOwner: {PermissionRead, PermissionWrite, PermissionDelete},
	RoleAdmin: {PermissionRead, PermissionWrite, PermissionDelete, PermissionAdmin},
}

// RequireRole creates middleware that requires specific roles
func RequireRole(apiCfg *handlers.ApiConfig, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context (set by auth middleware)
			userID, err := auth.GetUserIDFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get user details from database
			user, err := apiCfg.Queries.GetUserByID(r.Context(), userID)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// Check if user's role is in the allowed roles
			userRole := user.Role
			if slices.Contains(allowedRoles, userRole) {
				next.ServeHTTP(w, r)
				return
			}

			// Check role hierarchy - owner > admin > moderator > user
			if isRoleHigherOrEqual(userRole, allowedRoles) {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}

// RequirePermission creates middleware that requires specific permissions
func RequirePermission(apiCfg *handlers.ApiConfig, requiredPermission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context (set by auth middleware)
			userID, err := auth.GetUserIDFromContext(r.Context())
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get user details from database
			user, err := apiCfg.Queries.GetUserByID(r.Context(), userID)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// Check if user's role has the required permission
			userRole := user.Role
			permissions, exists := RolePermissions[userRole]
			if !exists {
				http.Error(w, "Invalid role", http.StatusForbidden)
				return
			}

			for _, permission := range permissions {
				if permission == requiredPermission {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}

// Helper function to check role hierarchy
func isRoleHigherOrEqual(userRole string, allowedRoles []string) bool {
	// Define role hierarchy
	roleHierarchy := map[string]int{
		RoleUser:  1,
		RoleMod:   2,
		RoleAdmin: 3,
		RoleOwner: 4,
	}

	userLevel, userExists := roleHierarchy[userRole]
	if !userExists {
		return false
	}

	// Check if user level meets any of the allowed role levels
	for _, allowedRole := range allowedRoles {
		if allowedLevel, exists := roleHierarchy[allowedRole]; exists {
			if userLevel >= allowedLevel {
				return true
			}
		}
	}

	return false
}

// Convenience middleware functions
func RequireAdminOrOwner(apiCfg *handlers.ApiConfig) func(http.Handler) http.Handler {
	return RequireRole(apiCfg, RoleAdmin, RoleOwner)
}

func RequireModeratorOrAbove(apiCfg *handlers.ApiConfig) func(http.Handler) http.Handler {
	return RequireRole(apiCfg, RoleMod, RoleAdmin, RoleOwner)
}
