package handler

import (
	"auth/internal/middleware"
	"auth/internal/models"
	"auth/internal/pkg"
	"auth/internal/services"
	"auth/internal/utils/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(as *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: as,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		validateErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(validateErrs))
		return
	}

	res, err := a.AuthService.Register(r.Context(), &req)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusCreated, res)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validate.Struct(req); err != nil {
		validateErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(validateErrs))
		return
	}

	res, err := a.AuthService.Login(r.Context(), &req)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusCreated, res)
}

func (a *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	resp, err := a.AuthService.Refresh(r.Context(), req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req models.LogoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	err := a.AuthService.Logout(req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.GeneralError(fmt.Errorf("Successfully logged out")))
}

func (a *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to logout from all devices")))
		return
	}

	err := a.AuthService.LogoutAll(req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, response.GeneralError(fmt.Errorf("Successfully logged out from all devices")))
}

func (a *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*pkg.Claims)
	if !ok {
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized")))
		return
	}

	response.WriteJSON(w, http.StatusOK, claims)
}
