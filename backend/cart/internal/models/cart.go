package models

import "time"

type Cart struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" validate:"required,gt=0"`
	VehicleId int       `json:"vehicle_id" validate:"required,gt=0"`
	Price     float64   `json:"price" validate:"required,gte=0"`
	Quantity  int       `json:"quantity" validate:"required,gte=1,lte=100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
