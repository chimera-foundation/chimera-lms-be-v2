package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
)


func AuthMiddleware(tokenProvider auth.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := tokenProvider.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			
			blacklisted, err := tokenProvider.IsBlacklisted(r.Context(), tokenString)
			if err != nil || blacklisted {
				http.Error(w, "Token has been revoked", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}