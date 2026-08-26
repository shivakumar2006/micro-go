package handler

import (
	"analytics/internal/kafka"
	"analytics/internal/service"
	"analytics/internal/utils/response"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	Service service.Service
}

func NewHandler(srv service.Service) *Handler {
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
		validateErr := err.(validator.ValidationErrors)
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
