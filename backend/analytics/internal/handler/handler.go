package handler

import (
	"analytics/internal/kafka"
	"analytics/internal/service"
	"analytics/internal/utils/response"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	Service *service.Service
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{
		Service: srv,
	}
}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (h *Handler) ProcessPaymentSuccess(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext()
	defer cancel()

	var req kafka.PaymentSuccessEvent

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if err := validate.Struct(req); err != nil {
		validateErr, ok := err.(validator.ValidationErrors)
		if !ok {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("validation failed")))
			return
		}
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request: %v", validateErr.Error())))
		return
	}

	if err := h.Service.ProcessPaymentSuccess(ctx, req); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "payment success event processed successfully",
	})
}

func (h *Handler) GetPaymentAnalytic(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext()
	defer cancel()

	data, err := h.Service.GetPaymentAnalytics(ctx)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) GetPaymentByPaymentID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext()
	defer cancel()

	idStr := chi.URLParam(r, "paymentID")
	if idStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("paymentID is required")))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid paymentID: %v", err)))
		return
	}

	data, err := h.Service.GetPaymentByPaymentID(ctx, id)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) GetPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext()
	defer cancel()

	idStr := chi.URLParam(r, "orderID")
	if idStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("orderID is required")))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid orderID: %v", err)))
		return
	}

	data, err := h.Service.GetPaymentByOrderID(ctx, id)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) GetPaymentByUserID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext()
	defer cancel()

	idStr := chi.URLParam(r, "userID")
	if idStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("userID is required")))
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid userID: %v", err)))
		return
	}

	data, err := h.Service.GetPaymentByUserID(ctx, id)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, data)
}
