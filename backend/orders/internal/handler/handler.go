package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"orders/internal/models"
	"orders/internal/services"
	"orders/internal/utils/response"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type OrderHandler struct {
	OrderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
	}
}

// POST    /orders
func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req models.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := validate.Struct(&req); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(validationErrs))
		return
	}

	order, err := o.OrderService.CreateOrder(ctx, &req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, order)
}

// GET     /orders/:id
// GET     /orders/user/:userId
// PATCH   /orders/:id/status
// PATCH   /orders/:id/cancel
// PATCH   /orders/:id/pay
