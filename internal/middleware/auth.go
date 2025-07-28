package middleware

import (
	"net/http"

	"github.com/omed0/go-hello-world/handlers"
	"github.com/omed0/go-hello-world/internal/auth"
)

// AuthMiddleware creates an authentication middleware that validates API keys
// and sets user context for protected routes
func AuthMiddleware(apiCfg *handlers.ApiConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract API key from Authorization header
			apiKey, err := auth.GetAPIKey(r.Header)
			if err != nil {
				handlers.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			// Validate API key and get user
			user, err := apiCfg.Queries.GetUserByAPIKey(r.Context(), apiKey)
			if err != nil {
				handlers.RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
				return
			}

			// Set user ID in context for downstream handlers
			ctx := auth.SetUserIDInContext(r.Context(), user.ID)
			r = r.WithContext(ctx)

			// Call next handler
			next.ServeHTTP(w, r)
		})
	}
}
