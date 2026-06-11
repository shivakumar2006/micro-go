package routes

import (
	db "auth/internal/db/storage"
	"auth/internal/model"
	"auth/internal/pkg"
	"auth/internal/utils/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo      db.Storage
	JwtManager    *pkg.JWTManager
	RefreshExpiry time.Duration
}

func NewAuthservice(userRepo db.Storage, jwtManager pkg.JWTManager, refreshExpiry time.Duration) *AuthService {
	return &AuthService{
		UserRepo:      userRepo,
		JwtManager:    &jwtManager,
		RefreshExpiry: refreshExpiry,
	}
}

func (a *AuthService) Register(storage db.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.RegisterRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		// request validation
		if err := validator.New().Struct(req); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidateErrors(validateErrs))
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to hashed password: %s", err)))
			return
		}

		user := model.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		if err := a.UserRepo.CreateUser(&user); err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to create user : %s", err)))
			return
		}

		tokenPair, err := a.JwtManager.GenerateTokenPair(strconv.FormatInt(user.ID, 10), user.Email)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to generate token pair : %s", err)))
		}

		refreshToken := model.RefreshToken{
			UserID:    user.ID,
			Token:     tokenPair.RefreshToken,
			ExpiresAt: time.Now().Add(a.RefreshExpiry),
		}

		if err := a.UserRepo.SaveRefreshToken(&refreshToken); err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to save refresh token : %s", err)))
			return
		}

		authResponse := buildAuthResponse(tokenPair, &user)

		response.WriteJSON(w, http.StatusCreated, authResponse)
	}
}

func buildAuthResponse(tokenPair *pkg.TokenPair, user *model.User) *model.AuthResponse {
	return &model.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User: model.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}
}
