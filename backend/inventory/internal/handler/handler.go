package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"inventory/internal/models"
	"inventory/internal/service"
	"inventory/internal/utils/response"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type InventoryHandler struct {
	Service *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		Service: svc,
	}
}

// POST   /api/v1/inventory
func (h *InventoryHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext()
	defer cancel()

	var req models.CreateInventoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if err := validate.Struct(req); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.ValidateErrors(validationErrs))
		return
	}

	err := h.Service.CreateTransaction(ctx, req.OrderID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Order inventory transaction created successfully",
	})
}

// GET    /api/v1/inventory/{id}
func (h *InventoryHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext()
	defer cancel()

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	transaction, err := h.Service.GetTransactionByID(ctx, int64(id))
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, transaction)
}

// GET    /api/v1/inventory/order/{orderId}
func (h *InventoryHandler) GetTransactionByOrderID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext()
	defer cancel()

	orderIDStr := chi.URLParam(r, "orderId")
	if orderIDStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	transaction, err := h.Service.GetTransactionByOrderID(ctx, int64(orderID))
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, transaction)
}

// PATCH  /api/v1/inventory/{id}

func (h *InventoryHandler) UpdateTransactionStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext()
	defer cancel()

	transactionIDStr := chi.URLParam(r, "id")
	if transactionIDStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	var req models.UpdateTransactionStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if err := validate.Struct(req); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.ValidateErrors(validationErrs))
		return
	}

	err = h.Service.UpdateTransactionStatus(ctx, int64(transactionID), req.Status)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Order inventory transaction updated successfully",
	})
}

// GET    /api/v1/inventory
func (h *InventoryHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext()
	defer cancel()

	txns, err := h.Service.GetTransactions(ctx)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, txns)
}

func newContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
