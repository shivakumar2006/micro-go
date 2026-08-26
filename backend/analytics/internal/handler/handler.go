package handler

import (
	"analytics/internal/kafka"
	"analytics/internal/service"
	"context"
	"encoding/json"
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
		
	}

	if err := validate.Struct()
}
