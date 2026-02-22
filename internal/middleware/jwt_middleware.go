package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kshitij-nehete/astro-report/internal/auth"
	"github.com/kshitij-nehete/astro-report/internal/handler"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func JWTMiddleware(jwtService *auth.JWTService) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				handler.WriteJSONError(w, http.StatusUnauthorized, "missing token")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 {
				handler.WriteJSONError(w, http.StatusUnauthorized, "invalid token format")
				return
			}

			token, err := jwtService.ValidateToken(parts[1])
			if err != nil || !token.Valid {
				handler.WriteJSONError(w, http.StatusUnauthorized, "invalid token")
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				handler.WriteJSONError(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				handler.WriteJSONError(w, http.StatusUnauthorized, "invalid token payload")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
