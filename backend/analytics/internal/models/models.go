package models

import "time"

type Analytics struct {
	ID        int64     `json:"id"`
	PaymentID int64     `json:"payment_id"`
	OrderID   int64     `json:"order_id"`
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	PaidAt    time.Time `json:"paid_at"`
	CreatedAt time.Time `json:"created_at"`
}
