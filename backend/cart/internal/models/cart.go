package models

import "time"

type Cart struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	VehicleId int       `json:"vehicle_id"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
