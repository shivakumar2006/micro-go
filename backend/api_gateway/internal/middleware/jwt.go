package middleware

import (
	"context"
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
	*jwt.RegisteredClaims
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
			http.Error(w, "authorization header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "authorization header format must be Bearer <token>", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := a.ParseToken(tokenString, a.AccessSecret)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		if claims.TokenType != "access" {
			http.Error(w, "access token is required", http.StatusUnauthorized)
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
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
