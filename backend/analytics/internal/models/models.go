package models

import "time"

type Analytics struct {
	ID        int64     `json:"id"`
	PaymentID int64     `json:"payment_id" validate:"required,gt=0"`
	OrderID   int64     `json:"order_id" validate:"required,gt=0"`
	UserID    int64     `json:"user_id" validate:"required,gt=0"`
	Email     string    `json:"email" validate:"required,email"`
	Status    string    `json:"status" validate:"required,eq=PAID"`
	PaidAt    time.Time `json:"paid_at" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}
