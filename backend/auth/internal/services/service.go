package services

import (
	"auth/internal/db/storage"
	"auth/internal/models"
	"auth/internal/pkg"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthRepo      storage.Storage
	JWTManager    *pkg.JWTManager
	RefreshExpiry time.Duration
}

func NewAuthService(repo storage.Storage, jwtmanager *pkg.JWTManager, refreshExpiry time.Duration) *AuthService {
	return &AuthService{
		AuthRepo:      repo,
		JWTManager:    jwtmanager,
		RefreshExpiry: refreshExpiry,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	exists, err := s.AuthRepo.FindUserByEmail(req.Email)
	if err == nil && exists != nil {
		return nil, errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to process password")
	}

	role := strings.TrimSpace(strings.ToLower(req.Role))

	switch role {
	case "customer", "admin":

	default:
		return nil, errors.New("invalid role")
	}

	roleID, err := s.AuthRepo.GetRoleByName(role)
	if err != nil {
		return nil, fmt.Errorf("failed to get default role : %w", err)
	}

	if role == "admin" {
		exists, err := s.AuthRepo.AdminExist()
		if err != nil {
			return nil, fmt.Errorf("failed to check admin exists : %w", err)
		}

		if exists {
			return nil, errors.New("admin already exists")
		}
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		RoleID:   roleID,
		RoleName: role,
	}

	if err := s.AuthRepo.CreateUser(user); err != nil {
		return nil, err
	}

	tokenPair, err := s.JWTManager.GenerateTokenPair(user.ID, user.Email, user.RoleName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair : %w", err)
	}

	// hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(tokenPair.RefreshToken), bcrypt.DefaultCost)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to process refresh token : %w", err)
	// }

	session := &models.UserSessions{
		UserID:           user.ID,
		RefreshTokenHash: HashRefreshToken(tokenPair.RefreshToken),
		ExpiresAt:        time.Now().Add(s.RefreshExpiry),
	}

	if err := s.AuthRepo.SaveRefreshToken(session); err != nil {
		return nil, err
	}

	return buildAuthResponse(tokenPair, user)
}

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (a *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := a.AuthRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email : %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	tokenPair, err := a.JWTManager.GenerateTokenPair(user.ID, user.Email, user.RoleName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair : %w", err)
	}

	// claims, err := a.JWTManager.ValidateAccessToken(tokenPair.AccessToken)
	// fmt.Println("SELF VALIDATION ERROR:", err)
	// fmt.Printf("SELF CLAIMS: %+v\n", claims)

	// hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(tokenPair.RefreshToken), bcrypt.DefaultCost)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to hashed refresh token : %w", err)
	// }

	session := &models.UserSessions{
		UserID:           user.ID,
		RefreshTokenHash: HashRefreshToken(tokenPair.RefreshToken),
		ExpiresAt:        time.Now().Add(a.RefreshExpiry),
	}

	if err := a.AuthRepo.SaveRefreshToken(session); err != nil {
		return nil, err
	}

	return buildAuthResponse(tokenPair, user)
}

func (a *AuthService) Refresh(ctx context.Context, req models.RefreshRequest) (*models.AuthResponse, error) {
	claims, err := a.JWTManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate refresh token : %w", err)
	}

	hashedToken := HashRefreshToken(req.RefreshToken)

	storedToken, err := a.AuthRepo.FindRefreshToken(hashedToken)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found, please login again : %w", err)
	}

	if time.Now().After(storedToken.ExpiresAt) {
		_ = a.AuthRepo.DeleteRefreshToken(hashedToken)
		return nil, fmt.Errorf("refresh token expired")
	}

	user, err := a.AuthRepo.FindUserById(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id : %w", err)
	}

	// delete old refresh token (rotation)
	_ = a.AuthRepo.DeleteRefreshToken(req.RefreshToken)

	// generate new token pair
	tokenPair, err := a.JWTManager.GenerateTokenPair(user.ID, user.Email, user.RoleName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair : %w", err)
	}

	// hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(tokenPair.RefreshToken), bcrypt.DefaultCost)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to process refresh token : %w", err)
	// }

	session := &models.UserSessions{
		UserID:           user.ID,
		RefreshTokenHash: HashRefreshToken(tokenPair.RefreshToken),
		ExpiresAt:        time.Now().Add(a.RefreshExpiry),
	}

	if err := a.AuthRepo.SaveRefreshToken(session); err != nil {
		return nil, err
	}

	return buildAuthResponse(tokenPair, user)
}

func (a *AuthService) Logout(req models.LogoutRequest) error {
	hashedToken := HashRefreshToken(req.RefreshToken)
	return a.AuthRepo.DeleteRefreshToken(hashedToken)
}

func (a *AuthService) LogoutAll(req models.RefreshRequest) error {
	hashedToken := HashRefreshToken(req.RefreshToken)

	token, err := a.AuthRepo.FindRefreshToken(hashedToken)
	if err != nil {
		return err
	}

	if err := a.AuthRepo.DeleteAllUserToken(token.UserID); err != nil {
		return err
	}

	return nil
}

func (a *AuthService) GetMe(userId int) (*models.UserResponse, error) {
	user, err := a.AuthRepo.FindUserById(userId)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        int64(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		RoleName:  user.RoleName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func buildAuthResponse(tokenPair *pkg.TokenPair, user *models.User) (*models.AuthResponse, error) {
	return &models.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User: models.UserResponse{
			ID:       int64(user.ID),
			Name:     user.Name,
			Email:    user.Email,
			RoleName: user.RoleName,
		},
	}, nil
}
