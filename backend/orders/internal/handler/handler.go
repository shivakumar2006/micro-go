package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"orders/internal/models"
	"orders/internal/services"
	"orders/internal/utils/response"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	ctx := r.Context()

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

	response.WriteJSON(w, http.StatusCreated, order)
}

func (o *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id")))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id")))
		return
	}

	order, err := o.OrderService.GetOrderByID(ctx, id)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, order)
}

func (o *OrderHandler) GetOrdersByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "user_Id")
	if idStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id")))
		return
	}

	userId, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id")))
		return
	}

	orders, err := o.OrderService.GetOrdersByUserID(ctx, userId)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, orders)
}

// PATCH   /orders/:id/status
func (o *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id")))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id")))
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	if err := validate.Struct(&req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(validateErr))
		return
	}

	if err := o.OrderService.UpdateOrderStatus(ctx, id, req.Status); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "order status updated successfully",
	})
}

func (o *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id")))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id")))
		return
	}

	if err := o.OrderService.CancelOrder(ctx, int64(id)); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "order cancel successfully",
	})
}

// PATCH   /orders/:id/pay
func (o *OrderHandler) MarkOrderPaid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id")))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id")))
		return
	}

	if err := o.OrderService.MarkOrderPaid(ctx, int64(id)); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "order paid successfully",
	})
}
