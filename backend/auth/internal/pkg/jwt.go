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
	UserID    int    `json:"user_id"`
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
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func NewJWTManager(as, rs string, ar, re time.Duration) (*JWTManager, error) {
	if as == "" {
		panic("access secret is required")
	}
	if rs == "" {
		panic("refresh secret is required")
	}

	if ar == 0 {
		panic("access expiry is required")
	}
	if re == 0 {
		panic("refresh expiry is required")
	}

	return &JWTManager{
		AccessSecret:  as,
		RefreshSecret: rs,
		AccessExpiry:  ar,
		RefreshExpiry: re,
	}, nil
}

func (j *JWTManager) GenerateTokenPair(userID int, email, role string) (*TokenPair, error) {
	accessToken, err := j.GenerateToken(userID, email, role, AccessToken, j.AccessSecret, j.AccessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token : %w", err)
	}

	refreshToken, err := j.GenerateToken(userID, email, role, RefreshToken, j.RefreshSecret, j.RefreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token : %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JWTManager) GenerateToken(userID int, email, role string, tokenType TokenType, secret string, expiry time.Duration) (string, error) {
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

	fmt.Println("Signing Secret:", secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (j *JWTManager) ValidateAccessToken(token string) (*Claims, error) {
	return j.ValidateToken(token, j.AccessSecret, AccessToken)
}

func (j *JWTManager) ValidateRefreshToken(token string) (*Claims, error) {
	return j.ValidateToken(token, j.RefreshSecret, RefreshToken)
}

func (j *JWTManager) ValidateToken(tokenString string, secret string, tokenType TokenType) (*Claims, error) {
	fmt.Println("===================================")
	fmt.Println("Incoming Token:", tokenString)
	fmt.Println("Using Secret:", secret)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			fmt.Println("Algorithm:", token.Method.Alg())
			return []byte(secret), nil
		},
	)

	if err != nil {
		fmt.Printf("PARSE ERROR: %#v\n", err)
		fmt.Println("ERROR STRING:", err.Error())
		return nil, fmt.Errorf("cannot parse token: %w", err)
	}

	fmt.Println("Token.Valid:", token.Valid)

	claims, ok := token.Claims.(*Claims)
	fmt.Printf("Claims: %+v\n", claims)

	if !ok {
		fmt.Println("Claims type assertion failed")
		return nil, errors.New("invalid claims")
	}

	if !token.Valid {
		fmt.Println("Token.Valid is false")
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// func (j *JWTManager) ValidateToken(tokenString string, secret string, tokenType TokenType) (*Claims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
// 		func(token *jwt.Token) (interface{}, error) {
// 			fmt.Println("Using Secret:", secret)
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, errors.New("unexpected signing method")
// 			}
// 			return ([]byte(secret)), nil
// 		},
// 	)

// 	if err != nil {
// 		return nil, fmt.Errorf("cannot parse the token : %w", err)
// 	}

// 	// extract claims from token
// 	claims, ok := token.Claims.(*Claims)
// 	if !ok || !token.Valid {
// 		return nil, errors.New("invalid token")
// 	}
// 	if claims.TokenType != string(tokenType) {
// 		return nil, errors.New("invalid token type")
// 	}

// 	if claims.ExpiresAt.Time.Before(time.Now()) {
// 		return nil, errors.New("token is expired")
// 	}

// 	return claims, nil
// }
