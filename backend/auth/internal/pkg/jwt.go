package pkg

import (
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
	Email     string `json:"gmail"`
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
	return &JWTManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (m *JWTManager) GenerateTokenPair(userID, email string) (*TokenPair, error) {
	accessToken, err := m.GenerateToken(userID, email, AccessToken, m.accessSecret, m.accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token %w", err)
	}

	refreshToken, err := m.GenerateToken(userID, email, RefreshToken, m.refreshSecret, m.refreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *JWTManager) GenerateToken(userID, email string, tokenType TokenType, secret string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Email:     email,
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
