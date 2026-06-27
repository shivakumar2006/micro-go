package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vehicles/internal/models"
	"vehicles/internal/service"
	"vehicles/internal/utils"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type VehicleHandler struct {
	VehicleService service.Service
}

// POST /vehicles
func (v *VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var req models.Vehicle

	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.GeneralError(fmt.Errorf("method not allowed")))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if err := validate.Struct(req); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		utils.WriteJSON(w, http.StatusBadRequest, utils.ValidateErrors(validationErrs))
		return
	}

	if err := v.VehicleService.CreateVehicle(&req); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.GeneralError(err))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Response{
		Status:  utils.StatusOK,
		Message: "Vehicle created successfully",
	})
}

// GET /vehicles

// GET /vehicles/:id

// PUT /vehicles/:id

// DELETE /vehicles/:id
