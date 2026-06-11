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

func NewAuthservice(userRepo db.Storage, jwtManager *pkg.JWTManager, refreshExpiry time.Duration) *AuthService {
	return &AuthService{
		UserRepo:      userRepo,
		JwtManager:    jwtManager,
		RefreshExpiry: refreshExpiry,
	}
}

func (a *AuthService) Register() http.HandlerFunc {
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
			return
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

func (a *AuthService) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		// request validation
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(validateErr))
			return
		}

		user, err := a.UserRepo.FindUserByEmail(req.Email)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid email credentials")))
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid password credentials")))
			return
		}

		tokenPair, err := a.JwtManager.GenerateTokenPair(strconv.FormatInt(user.ID, 10), user.Email)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to generate token pair : %s", err)))
			return
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

		authResponse := buildAuthResponse(tokenPair, user)
		response.WriteJSON(w, http.StatusOK, authResponse)
	}
}

func (a *AuthService) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.RefreshRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidateErrors(validateErr))
			return
		}

		// validate jwt + expiry
		claims, err := a.JwtManager.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid refresh token : %s", err)))
			return
		}

		token, err := a.UserRepo.FindRefreshToken(req.RefreshToken)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid refresh token : %s", err)))
			return
		}

		// check if the token is expired or not in db
		if time.Now().After(token.ExpiresAt) {
			_ = a.UserRepo.DeleteRefreshToken(req.RefreshToken)
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("refresh token expired")))
			return
		}

		// get user from the db
		userID, err := strconv.ParseInt(claims.UserID, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid token claims")))
			return
		}

		user, err := a.UserRepo.FindUserById(userID)
		if err != nil {
			response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("user no longer exist")))
			return
		}

		// delete old refresh token
		_ = a.UserRepo.DeleteRefreshToken(req.RefreshToken)

		// generate new token pair
		tokenPair, err := a.JwtManager.GenerateTokenPair(strconv.FormatInt(user.ID, 10), user.Email)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to generate token pair : %s", err)))
			return
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

		authResponse := buildAuthResponse(tokenPair, user)

		response.WriteJSON(w, http.StatusOK, authResponse)

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
