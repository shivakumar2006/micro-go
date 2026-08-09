package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"payment/internal/models"
	"payment/internal/service"
	"payment/internal/utils/response"
	"time"
)

type PaymentHandler struct {
	Service *service.PaymentService
}

func NewPaymentHandler(s *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{Service: s}
}

// POST /payments/checkout
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req models.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body : %w", err)))
		return
	}

	payment, err := h.Service.CreatePayment(ctx, &req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"id":  payment.ID,
		"url": payment.URL,
	})
}

// GET /payments/{id}

// GET /payments/order/{orderId}

// POST /payments/webhook
