package middleware

import (
	"context"
	"net/http"

	"github.com/denga/go-real-world-example/internal/auth"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// UserEmailKey is the key used to store the user email in the request context
const UserEmailKey contextKey = "userEmail"

// Auth is middleware that validates JWT tokens and adds the user email to the request context
func Auth(config auth.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for certain endpoints
			if r.URL.Path == "/api/users" && r.Method == http.MethodPost {
				// Registration endpoint
				next.ServeHTTP(w, r)
				return
			}
			if r.URL.Path == "/api/users/login" && r.Method == http.MethodPost {
				// Login endpoint
				next.ServeHTTP(w, r)
				return
			}
			if r.URL.Path == "/api/tags" && r.Method == http.MethodGet {
				// Tags endpoint (public)
				next.ServeHTTP(w, r)
				return
			}
			if r.URL.Path == "/api/articles" && r.Method == http.MethodGet {
				// Articles listing (public)
				next.ServeHTTP(w, r)
				return
			}
			if r.URL.Path == "/openapi.yml" {
				// OpenAPI spec (public)
				next.ServeHTTP(w, r)
				return
			}

			// Extract token from request
			tokenString, err := auth.ExtractTokenFromRequest(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate token
			email, err := auth.ValidateToken(tokenString, config)
			if err != nil {
				if err == auth.ErrExpiredToken {
					http.Error(w, "Token expired", http.StatusUnauthorized)
				} else {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
				}
				return
			}

			// Add user email to context
			ctx := context.WithValue(r.Context(), UserEmailKey, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserEmail extracts the user email from the request context
func GetUserEmail(r *http.Request) (string, bool) {
	email, ok := r.Context().Value(UserEmailKey).(string)
	return email, ok
}
