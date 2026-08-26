package kafka

import "time"

type PaymentSuccessEvent struct {
	OrderID   int64     `json:"order_id" validate:"required,gt=0"`
	PaymentID int64     `json:"payment_id" validate:"required,gt=0"`
	UserID    int64     `json:"user_id" validate:"required,gt=0"`
	Email     string    `json:"email" validate:"required,email"`
	Status    string    `json:"status" validate:"required,eq=PAID"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}
