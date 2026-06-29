package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"vehicles/internal/models"
	"vehicles/internal/service"
	"vehicles/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type VehicleHandler struct {
	VehicleService *service.Service
}

func NewVehicleHandler(vehicleService *service.Service) *VehicleHandler {
	return &VehicleHandler{
		VehicleService: vehicleService,
	}
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
func (v *VehicleHandler) GetAllVehicles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.GeneralError(fmt.Errorf("method not allowed")))
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	query := models.VehicleQuery{
		Page:     page,
		Limit:    limit,
		Search:   r.URL.Query().Get("search"),
		Type:     r.URL.Query().Get("type"),
		Category: r.URL.Query().Get("category"),
		SortBy:   r.URL.Query().Get("sort_by"),
		Order:    r.URL.Query().Get("order"),
	}

	response, err := v.VehicleService.GetAllVehicles(query)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.GeneralError(err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GET /vehicles/:id
func (v *VehicleHandler) GetVehicleById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.GeneralError(fmt.Errorf("method not allowed")))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.GeneralError(fmt.Errorf("invalid vehicle id")))
		return
	}

	vehicle, err := v.VehicleService.GetVehicleById(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.GeneralError(err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, vehicle)
}

// PUT /vehicles/:id
func (v *VehicleHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.GeneralError(fmt.Errorf("method not allowed")))
		return
	}

	var req models.Vehicle

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if err := validate.Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		utils.WriteJSON(w, http.StatusBadRequest, utils.ValidateErrors(validateErr))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.GeneralError(fmt.Errorf("invalid vehicle id")))
		return
	}

	req.Id = id
	if err := v.VehicleService.UpdateVehicle(&req); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.GeneralError(err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Status:  utils.StatusOK,
		Message: "Vehicle updated successfully",
	})
}

// DELETE /vehicles/:id
func (v *VehicleHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, utils.GeneralError(fmt.Errorf("method not allowed")))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.GeneralError(fmt.Errorf("invalid vehicle id")))
		return
	}

	if err := v.VehicleService.DeleteVehicle(id); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.GeneralError(err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Status:  utils.StatusOK,
		Message: "Vehicle deleted successfully",
	})
}
