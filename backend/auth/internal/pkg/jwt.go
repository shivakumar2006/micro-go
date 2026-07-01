package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type JWTManager struct {
	accessSecret  string
	refreshSecret string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTManager(accessSecret, refreshSecret string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	if accessSecret == "" {
		panic("access secret is required")
	}

	if refreshSecret == "" {
		panic("refresh secret is required")
	}

	return &JWTManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (m *JWTManager) GenerateTokenPair(userID, email, role string) (*TokenPair, error) {
	accessToken, err := m.GenerateToken(userID, email, role, AccessToken, m.accessSecret, m.accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token %w", err)
	}

	refreshToken, err := m.GenerateToken(userID, email, role, RefreshToken, m.refreshSecret, m.refreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *JWTManager) GenerateToken(userID, email, role string, tokenType TokenType, secret string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "micro-go",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (m *JWTManager) ValidateAccessToken(token string) (*Claims, error) {
	return m.ValidateToken(token, m.accessSecret, AccessToken)
}

func (m *JWTManager) ValidateRefreshToken(token string) (*Claims, error) {
	return m.ValidateToken(token, m.refreshSecret, RefreshToken)
}

func (m *JWTManager) ValidateToken(tokenString string, secret string, tokenType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing token method")
			}
			return ([]byte(secret)), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("cannot parse the token %w", err)
	}

	// extract claims from token
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.TokenType != string(tokenType) {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
