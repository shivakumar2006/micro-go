package models

import "time"

type Notification struct {
	ID        int       `json:"id" db:"id"`
	OrderID   int       `json:"order_id" db:"order_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	Type      string    `json:"type" db:"type"`
	Status    string    `json:"status" db:"status"`
	Provider  string    `json:"provider" db:"provider"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	SentAt    time.Time `json:"sent_at" db:"sent_at"`
}
