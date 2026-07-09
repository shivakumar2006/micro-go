package handler

import (
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
	var cart models.Cart

	if r.Method != http.MethodPost {
		response.WriteJSON(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("method not allowed")))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request")))
		return
	}

	if err := validate.Struct(cart); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.ValidateErrors(validationErrs))
		return
	}

	err := c.Service.AddToCart(&cart)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to add to cart: %v", err)))
		return
	}

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "Item added to cart successfully",
	})
}

func (c *CartHandler) UpdateCartQuantity(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart

	if r.Method != http.MethodPut {
		response.WriteJSON(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("method not allowed")))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if err := validate.Struct(cart); err != nil {
		validateErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.ValidateErrors(validateErrs))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid cart id")))
		return
	}

	cart.ID = int(id)

	_, err = c.Service.UpdateCartQuantity(int(id), cart.Quantity)
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
	var cart models.Cart

	if r.Method != http.MethodGet {
		response.WriteJSON(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("method not allowed")))
		return
	}

	cartItems, err := c.Service.GetUserCart(cart.UserID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, cartItems)
}

func (c *CartHandler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.WriteJSON(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("method not allowed")))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid cart id")))
		return
	}

	quantity, err := c.Service.DeleteCartItem(int(id))
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
	if r.Method != http.MethodDelete {
		response.WriteJSON(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("method not allowed")))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
		return
	}

	deletedID, err := c.Service.DeleteCart(int(id))
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
	if r.Method != http.MethodGet {
		response.WriteJSON(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("invalid id")))
		return
	}

	var cart models.Cart

	cartTotal, err := c.Service.GetCartTotal(cart.UserID)
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
	if r.Method != http.MethodGet {
		response.WriteJSON(w, http.StatusMethodNotAllowed, response.GeneralError(fmt.Errorf("invalid id")))
		return
	}

	var cart models.Cart

	countItems, err := c.Service.CountItems(cart.UserID)
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
