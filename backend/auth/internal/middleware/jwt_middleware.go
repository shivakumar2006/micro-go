package middleware

import (
	"auth/internal/pkg"
	"auth/internal/utils/response"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type ContextKey string

const UserContextKey ContextKey = "user"

type AuthMiddleware struct {
	JwtManager *pkg.JWTManager
}

func NewAuthMiddleware(jwtManager pkg.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		JwtManager: &jwtManager,
	}
}

func (a *AuthMiddleware) Authnticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("authorization header is required")))
			return
		}

		// check the format of a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("format must be : Bearer <token>")))
			return
		}

		tokenString := parts[1]

		// validate the token
		claims, err := a.JwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid token or expired token : %s", err)))
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
