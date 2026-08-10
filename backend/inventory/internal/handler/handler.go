package handler

import (
	"encoding/json"
	"fmt"
	"inventory/internal/models"
	"inventory/internal/service"
	"inventory/internal/utils/response"
	"net/http"

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
}

// GET    /api/v1/inventory/{id}

// GET    /api/v1/inventory/order/{orderId}

// PATCH  /api/v1/inventory/{id}

// GET    /api/v1/inventory
