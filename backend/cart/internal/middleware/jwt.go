package middleware

import (
	"cart/internal/utils/response"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	UserIDKey      contextKey = "UserID"
	UserEmailIDKey contextKey = "UserEmail"
	RoleKey        contextKey = "Role"
)

type Claims struct {
	UserID    int    `json:"user_id"`
	UserEmail string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	AccessSecret  string `json:"access-secret"`
	RefreshSecret string `json:"refresh-token"`
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
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("auth header is missing")))
			return
		}

		// check format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("auth header is not in valid format")))
			return
		}

		tokenString := parts[1]

		claims, err := a.ParseToken(tokenString, a.AccessSecret)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid claims")))
			return
		}

		if claims.TokenType != "access" {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("access token is required")))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailIDKey, claims.UserEmail)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)

		r.Header.Set("X-User-ID", strconv.Itoa(claims.UserID))
		r.Header.Set("X-User-Email", claims.UserEmail)
		r.Header.Set("X-Role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) ParseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}

func (a *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("role not found")))
				return
			}

			if userRole != role {
				response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("access denied")))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
