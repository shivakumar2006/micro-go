package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"orders/internal/utils/response"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const InternalServiceKey contextKey = "internal_service"

const (
	UserIDKey    contextKey = "UserId"
	UserEmailKey contextKey = "email"
	RoleKey      contextKey = "role"
)

type Claims struct {
	UserID    int    `json:"user_id"`
	UserEmail string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	AccessSecret       string `json:"access-secret"`
	RefreshSecret      string `json:"refresh-secret"`
	InternalServiceKey string `json:"internal-service-key"`
}

func NewAuthMiddleware(as, rs, isk string) *AuthMiddleware {
	return &AuthMiddleware{
		AccessSecret:       as,
		RefreshSecret:      rs,
		InternalServiceKey: isk,
	}
}

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Internal service request
		if authHeader == "Internal "+a.InternalServiceKey {
			ctx := context.WithValue(
				r.Context(),
				InternalServiceKey,
				true,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if authHeader == "" {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(errors.New("authorization header is missing")))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(errors.New("invalid auth header format")))
			return
		}

		tokenString := parts[1]

		claim, err := a.parseToken(tokenString, a.AccessSecret)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("failed to parse token : %w", err)))
			return
		}

		if claim.TokenType != "access" {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("access token is required")))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claim.UserID)
		ctx = context.WithValue(r.Context(), UserEmailKey, claim.UserEmail)
		ctx = context.WithValue(r.Context(), RoleKey, claim.Role)

		r.Header.Set("X-User-ID", strconv.Itoa(claim.UserID))
		r.Header.Set("X-User-Email", claim.UserEmail)
		r.Header.Set("X-Role", claim.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) parseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method : %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})

	if err != nil {
		return nil, fmt.Errorf("token is invalid : %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claim, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims : %w", err)
	}

	return claim, nil
}

func (a *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if internal, ok := r.Context().Value(InternalServiceKey).(bool); ok && internal {
				next.ServeHTTP(w, r)
				return
			}

			userRole, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(errors.New("user role not found in context")))
				return
			}

			if userRole != role {
				response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("required role : %s, but got : %s", role, userRole)))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
