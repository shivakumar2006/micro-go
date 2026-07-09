package handler

import (
	"cart/internal/middleware"
	"cart/internal/models"
	"cart/internal/service"
	"cart/internal/utils/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CartHandler struct {
	Service *service.CartService
}

func NewCartHandler(service *service.CartService) *CartHandler {
	return &CartHandler{
		Service: service,
	}
}

func (c *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req models.AddToCartRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request")))
		return
	}

	if err := validate.Struct(req); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.ValidateErrors(validationErrs))
		return
	}

	userId, err := getUserID(r)
	if err != nil {
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	cart := models.Cart{
		UserID:    userId,
		VehicleId: req.VehicleID,
		Quantity:  req.Quantity,
	}

	err = c.Service.AddToCart(&cart)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to add to cart: %v", err)))
		return
	}

	response.WriteJSON(w, http.StatusCreated, &response.Response{
		Status:  response.StatusOK,
		Message: "Item added to cart successfully",
	})
}

func (c *CartHandler) UpdateCartQuantity(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateQuantityRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if err := validate.Struct(req); err != nil {
		validateErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.ValidateErrors(validateErrs))
		return
	}

	idStr := chi.URLParam(r, "id")
	cartID, _ := strconv.Atoi(idStr)

	userId, err := getUserID(r)
	if err != nil {
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
	}

	_, err = c.Service.UpdateCartQuantity(userId, cartID, req.Quantity)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "Cart quantity updated successfully",
	})
}

func (c *CartHandler) GetUserCart(w http.ResponseWriter, r *http.Request) {

	userID, err := getUserID(r)
	if err != nil {
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	cartItems, err := c.Service.GetUserCart(userID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, cartItems)
}

func (c *CartHandler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	cartID, _ := strconv.Atoi(idStr)

	userId, err := getUserID(r)
	if err != nil {
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	quantity, err := c.Service.DeleteCartItem(userId, cartID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "Cart item deleted successfully",
		Data:    quantity,
	})
}

func (c *CartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {

	userId, err := getUserID(r)
	if err != nil {
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	deletedID, err := c.Service.DeleteCart(userId)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to delete cart: %v", err)))
		return
	}

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "Cart deleted successfully",
		Data:    deletedID,
	})
}

func (c *CartHandler) GetCartTotal(w http.ResponseWriter, r *http.Request) {

	userId, err := getUserID(r)
	if err != nil {
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	cartTotal, err := c.Service.GetCartTotal(userId)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to get all cart total : %v", err)))
		return
	}

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "all cart total retrive successfully",
		Data:    cartTotal,
	})
}

func (c *CartHandler) CountItems(w http.ResponseWriter, r *http.Request) {

	userId, err := getUserID(r)
	if err != nil {
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	countItems, err := c.Service.CountItems(userId)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to count items in cart : %v", err)))
		return
	}

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "all cart total retrive successfully",
		Data:    countItems,
	})
}

func getUserID(r *http.Request) (int, error) {
	userIdStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		return 0, fmt.Errorf("unauthorized")
	}

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		return 0, fmt.Errorf("invalid user id")
	}

	return userID, nil
}
