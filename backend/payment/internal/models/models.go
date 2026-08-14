package models

import "time"

const (
	StatusPending   = "PENDING"
	StatusPaid      = "PAID"
	StatusCancelled = "CANCELLED"
	StatusDelivered = "DELIVERED"
)

type Payment struct {
	ID              int       `json:"id"`
	OrderID         int       `json:"order_id"`
	UserID          int       `json:"user_id"`
	Email           string    `json:"email"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Provider        string    `json:"provider"`
	PaymentIntentID *string   `json:"payment_intent_id"`
	StripeSessionID string    `json:"stripe_session_id"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type OrderResponse struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

type StripeCheckoutResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type CreatePaymentRequest struct {
	OrderID int `json:"order_id"`
}

type StripeSession struct {
}
