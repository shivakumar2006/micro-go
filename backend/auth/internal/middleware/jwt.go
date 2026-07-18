package middleware

import (
	"auth/internal/pkg"
	"context"
	"net/http"
	"strings"
)

type ContextKey string

const UserContextKey ContextKey = "user"

type AuthMiddleware struct {
	JwtManager *pkg.JWTManager
}

func NewAuthMiddleware(jwt *pkg.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		JwtManager: jwt,
	}
}

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorized header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "format must be: Bearer <token>", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := a.JwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
