package kafka

import "time"

type PaymentSuccessEvent struct {
	OrderID   int64     `json:"order_id"`
	PaymentID int64     `json:"payment_id"`
	UserID    int64     `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
