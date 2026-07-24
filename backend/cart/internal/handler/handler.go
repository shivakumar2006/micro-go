package handler

import (
	"cart/internal/middleware"
	"cart/internal/models"
	"cart/internal/service"
	"cart/internal/utils/response"
	"encoding/json"
	"fmt"
	"log/slog"
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
	slog.Info("Add to cart request received", slog.String("method", r.Method), slog.String("url", r.URL.Path))

	var req models.AddToCartRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		slog.Warn("invalid request body", slog.String("error", err.Error()))

		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request")))
		return
	}

	if err := validate.Struct(req); err != nil {

		slog.Warn("Validation failed", slog.Any("error", err))

		validationErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.ValidateErrors(validationErrs))
		return
	}

	userId, err := getUserID(r)
	if err != nil {
		slog.Warn("Unauthorized request", slog.String("error", err.Error()))
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
		slog.Error("Failed to add item to cart", slog.Int("user_id", userId), slog.Int("vehicle_id", req.VehicleID), slog.String("error", err.Error()))

		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to add to cart: %v", err)))
		return
	}

	slog.Info("Item added to cart successfully", slog.Int("user_id", userId), slog.Int("vehicle_id", req.VehicleID))

	response.WriteJSON(w, http.StatusCreated, &response.Response{
		Status:  response.StatusOK,
		Message: "Item added to cart successfully",
	})
}

func (c *CartHandler) UpdateCartQuantity(w http.ResponseWriter, r *http.Request) {
	slog.Info("Update cart quantity request received", slog.String("method", r.Method), slog.String("url", r.URL.Path))

	var req models.UpdateQuantityRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("invalid request body", slog.String("error", err.Error()))

		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body")))
		return
	}

	if err := validate.Struct(req); err != nil {
		slog.Warn("Validation failed", slog.Any("error", err))

		validateErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.ValidateErrors(validateErrs))
		return
	}

	idStr := chi.URLParam(r, "id")
	cartID, _ := strconv.Atoi(idStr)

	userId, err := getUserID(r)
	if err != nil {
		slog.Warn("Unauthorized request", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	_, err = c.Service.UpdateCartQuantity(userId, cartID, req.Quantity)
	if err != nil {
		slog.Error("Failed to update cart quantity", slog.Int("user_id", userId), slog.Int("cart_id", cartID), slog.Int("quantity", req.Quantity), slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	slog.Info("Cart quantity updated successfully", slog.Int("user_id", userId), slog.Int("cart_id", cartID), slog.Int("quantity", req.Quantity))

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "Cart quantity updated successfully",
	})
}

func (c *CartHandler) GetUserCart(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get user cart request received", slog.String("method", r.Method), slog.String("url", r.URL.Path))

	userID, err := getUserID(r)
	if err != nil {
		slog.Warn("Unauthorized request", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	cartItems, err := c.Service.GetUserCart(userID)
	if err != nil {
		slog.Error("Failed to get user cart", slog.Int("user_id", userID), slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	slog.Info("User cart retrieved successfully", slog.Int("user_id", userID), slog.Int("cart_item_count", len(cartItems)))

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "User cart retrieved successfully",
		Data:    cartItems,
	})
}

func (c *CartHandler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	slog.Info("Delete cart item request received", slog.String("method", r.Method), slog.String("url", r.URL.Path))

	idStr := chi.URLParam(r, "id")
	cartID, _ := strconv.Atoi(idStr)

	userId, err := getUserID(r)
	if err != nil {
		slog.Warn("Unauthorized request", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	quantity, err := c.Service.DeleteCartItem(userId, cartID)
	if err != nil {
		slog.Error("Failed to delete cart item", slog.Int("user_id", userId), slog.Int("cart_id", cartID), slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
		return
	}

	slog.Info("Cart item deleted successfully", slog.Int("user_id", userId), slog.Int("cart_id", cartID), slog.Int("quantity", quantity))

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "Cart item deleted successfully",
		Data:    quantity,
	})
}

func (c *CartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	slog.Info("Delete cart request received", slog.String("method", r.Method), slog.String("url", r.URL.Path))

	userId, err := getUserID(r)
	if err != nil {
		slog.Warn("Unauthorized request", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	deletedID, err := c.Service.DeleteCart(userId)
	if err != nil {
		slog.Error("Failed to delete cart", slog.Int("user_id", userId), slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to delete cart: %v", err)))
		return
	}

	slog.Info("Cart deleted successfully", slog.Int("user_id", userId), slog.Int("deleted_id", deletedID))

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "Cart deleted successfully",
		Data:    deletedID,
	})
}

func (c *CartHandler) GetCartTotal(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get cart total request received", slog.String("method", r.Method), slog.String("url", r.URL.Path))

	userId, err := getUserID(r)
	if err != nil {
		slog.Warn("Unauthorized request", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	cartTotal, err := c.Service.GetCartTotal(userId)
	if err != nil {
		slog.Error("Failed to get cart total", slog.Int("user_id", userId), slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to get all cart total : %v", err)))
		return
	}

	slog.Info("Cart total retrieved successfully", slog.Int("user_id", userId), slog.Float64("cart_total", cartTotal))

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "all cart total retrive successfully",
		Data:    cartTotal,
	})
}

func (c *CartHandler) CountItems(w http.ResponseWriter, r *http.Request) {
	slog.Info("Count items request received", slog.String("method", r.Method), slog.String("url", r.URL.Path))

	userId, err := getUserID(r)
	if err != nil {
		slog.Warn("Unauthorized request", slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(err))
		return
	}

	countItems, err := c.Service.CountItems(userId)
	if err != nil {
		slog.Error("Failed to count items in cart", slog.Int("user_id", userId), slog.String("error", err.Error()))
		response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to count items in cart : %v", err)))
		return
	}

	slog.Info("Items counted successfully", slog.Int("user_id", userId), slog.Int("count", countItems))

	response.WriteJSON(w, http.StatusOK, &response.Response{
		Status:  response.StatusOK,
		Message: "all cart total retrive successfully",
		Data:    countItems,
	})
}

func getUserID(r *http.Request) (int, error) {
	slog.Info("Get user id request received", slog.String("method", r.Method), slog.String("url", r.URL.Path))

	userIdStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		slog.Warn("Unauthorized request", slog.String("error", "unauthorized"))
		return 0, fmt.Errorf("unauthorized")
	}

	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		slog.Warn("Invalid user id", slog.String("error", err.Error()))
		return 0, fmt.Errorf("invalid user id")
	}

	slog.Info("User id retrieved successfully", slog.Int("user_id", userID))

	return userID, nil
}
