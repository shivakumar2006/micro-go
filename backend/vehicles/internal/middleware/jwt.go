package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"vehicles/internal/utils"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	UserIDKey      contextKey = "UserID"
	UserEmailIDKey contextKey = "UserEmail"
	RoleKey        contextKey = "Role"
)

type Claims struct {
	UserId    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	JwtSecret     string
	RefreshSecret string
}

func NewAuthMiddleware(jwtSecret, refreshSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		JwtSecret:     jwtSecret,
		RefreshSecret: refreshSecret,
	}
}

func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.GeneralError(fmt.Errorf("authorization header is required")))
			return
		}

		// check format has bearer or not
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.GeneralError(fmt.Errorf("format must be: bearer <token>")))
			return
		}

		tokenString := parts[1]

		claims, err := parseToken(tokenString, a.JwtSecret)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.GeneralError(fmt.Errorf("invalid or expired token")))
			return
		}

		if claims.TokenType != "access" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.GeneralError(fmt.Errorf("invalid token type")))
			return
		}

		// add claims into request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserId)
		ctx = context.WithValue(ctx, UserEmailIDKey, claims.UserEmail)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)

		r.Header.Set("X-User-ID", claims.UserId)
		r.Header.Set("X-User-Email", claims.UserEmail)
		r.Header.Set("X-User-Role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func (a *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				utils.WriteJSON(w, http.StatusUnauthorized, utils.GeneralError(fmt.Errorf("role not found")))
				return
			}

			if userRole != role {
				utils.WriteJSON(w, http.StatusUnauthorized, utils.GeneralError(fmt.Errorf("access denied")))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
