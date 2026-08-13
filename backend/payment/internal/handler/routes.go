package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"payment/internal/models"
	"payment/internal/service"
	"payment/internal/utils/response"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

	authHeader := r.Header.Get("Authorization")

	var req models.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body : %w", err)))
		return
	}

	payment, err := h.Service.CreatePayment(ctx, &req, authHeader)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}
	slog.Info(
		"Creating payment",
		slog.Int("order_id", req.OrderID),
	)

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"id":  payment.ID,
		"url": payment.URL,
	})
}

// GET /payments/{id}
func (h *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	paymentID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid payment id : %w", err)))
		return
	}

	payment, err := h.Service.GetPaymentByID(ctx, paymentID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, payment)
}

// GET /payments/order/{orderId}
func (h *PaymentHandler) GetPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderID, err := strconv.Atoi(chi.URLParam(r, "orderId"))
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid order id : %w", err)))
		return
	}

	payment, err := h.Service.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, payment)
}

// POST /payments/webhook
func (h *PaymentHandler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body : %w", err)))
		return
	}
	defer r.Body.Close()

	// get stripe signature
	signature := r.Header.Get("Stripe-Signature")
	if signature == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing stripe signature : %w", err)))
		return
	}

	if err := h.Service.HandleWebhook(ctx, payload, signature); err != nil {
		slog.Error("webhook processing failed", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Webhook processed successfully",
	})
}
