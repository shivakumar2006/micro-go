package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payment/internal/models"
	"payment/internal/service"
	"payment/internal/utils/response"
)

type PaymentHandler struct {
	Service *service.PaymentService
}

func NewPaymentHandler(s *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{Service: s}
}

// POST /payments/checkout
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req models.Payment
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body : %w", err)))
		return
	}

}

// GET /payments/{id}

// GET /payments/order/{orderId}

// POST /payments/webhook
