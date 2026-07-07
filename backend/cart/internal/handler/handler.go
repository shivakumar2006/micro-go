package handler

import (
	"cart/internal/models"
	"cart/internal/service"
	"cart/internal/utils/response"
	"encoding/json"
	"fmt"
	"net/http"
)

type CartHandler struct {
	Service *service.CartService
}

func NewCartHandler(service *service.CartService) *CartHandler {
	return &CartHandler{
		Service: service,
	}
}

func (c *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var cart *models.Cart

	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request")))
		return
	}

	err := c.Service.AddToCart(cart)
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
	var cart *models.Cart

	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	_, err := c.Service.UpdateCartQuantity(cart.ID, cart.Quantity)
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

	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	_, err := c.Service.GetUserCart(cart.UserID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "User cart fetched successfully",
	})
}
