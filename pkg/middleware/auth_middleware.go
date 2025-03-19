package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/adilm/money-pulse/pkg/discovery"
	apperrors "github.com/adilm/money-pulse/pkg/errors"
)

// AuthMiddleware validates tokens by calling the auth service
func AuthMiddleware(consulClient *discovery.ConsulClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				apperrors.RespondWithError(w, apperrors.NewUnauthorizedError("Missing or invalid authorization header"))
				return
			}

			// Find auth service
			authServiceURL, err := consulClient.DiscoverService("auth-service")
			if err != nil {
				apperrors.RespondWithError(w, apperrors.NewInternalError(fmt.Errorf("auth service discovery failed: %w", err)))
				return
			}

			// Create new request to auth service
			req, err := http.NewRequest("GET", authServiceURL+"/auth/validate", nil)
			if err != nil {
				apperrors.RespondWithError(w, apperrors.NewInternalError(err))
				return
			}

			// Pass the token
			req.Header.Set("Authorization", authHeader)

			// Call auth service
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				apperrors.RespondWithError(w, apperrors.NewInternalError(fmt.Errorf("auth service unavailable: %w", err)))
				return
			}
			defer resp.Body.Close()

			// Check auth result
			if resp.StatusCode != http.StatusOK {
				apperrors.RespondWithError(w, apperrors.NewUnauthorizedError("Invalid token"))
				return
			}

			// Extract user info
			var userInfo map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
				apperrors.RespondWithError(w, apperrors.NewInternalError(err))
				return
			}

			// Add user info to context for later use
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", userInfo["user_id"])
			ctx = context.WithValue(ctx, "email", userInfo["email"])
			ctx = context.WithValue(ctx, "role", userInfo["role"])

			// Call next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole middleware checks if user has the required role
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("role").(string)
			if !ok || userRole != role {
				apperrors.RespondWithError(w, apperrors.NewForbiddenError("Insufficient permissions"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
