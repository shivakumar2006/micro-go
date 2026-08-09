package middleware

import (
	"context"
	"fmt"
	"net/http"
	"payment/internal/utils/response"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	UserIDKey    contextKey = "UserID"
	UserEmailKey contextKey = "UserEmail"
	RoleKey      contextKey = "Role"
)

const (
	Bearer = "Bearer"
)

type Claims struct {
	UserID    int    `json:"user_id"`
	UserEmail string `json:"user_email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	AccessSecret  string
	RefreshSecret string
}

func NewAuthMiddleware(as, rs string) *AuthMiddleware {
	return &AuthMiddleware{
		AccessSecret:  as,
		RefreshSecret: rs,
	}
}

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("auth header is required")))
			return
		}

		parts := strings.SplitN(authHeader, "", 2)
		if len(parts) != 2 || parts[0] != Bearer {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("auth header is not in a valid format")))
			return
		}

		tokenString := parts[1]

		claims, err := a.parseToken(tokenString, a.AccessSecret)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("auth token is not valid")))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), UserIDKey, claims.UserID))
		r = r.WithContext(context.WithValue(r.Context(), UserEmailKey, claims.UserEmail))
		r = r.WithContext(context.WithValue(r.Context(), RoleKey, claims.Role))

		next.ServeHTTP(w, r)
	})
}

func (a *AuthMiddleware) parseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (a *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			useRole, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("role is not valid")))
				return
			}
			if useRole != role {
				response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("user is not authorized")))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
