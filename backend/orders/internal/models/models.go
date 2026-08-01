package models

import "time"

const (
	OrderStatusPending   = "PENDING"
	OrderStatusPaid      = "PAID"
	OrderStatusCancelled = "CANCELLED"
	OrderStatusDelivered = "DELIVERED"
)

type Order struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	VehicleID int       `json:"vehicle_id"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
	UserID int64             `json:"user_id" validate:"required,min=1"`
	Items  []CreateOrderItem `json:"items" validate:"required,min=1"`
}

type CreateOrderItem struct {
	VehicleID int     `json:"vehicle_id" validate:"required,min=1"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=PENDING PAID CANCELLED DELIVERED"`
}

// cart service
type CartResponse struct {
	Status  int         `json:"cart"`
	Message string      `json:"message"`
	Data    []CartItems `json:"data"`
}

type CartItems struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	VehicleID int     `json:"vehicle_id"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}
