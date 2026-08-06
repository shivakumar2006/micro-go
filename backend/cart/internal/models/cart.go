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

type UpdateQuantityRequest struct {
	Quantity int `json:"quantity" validate:"required,gte=1"`
}

type AddToCartRequest struct {
	VehicleID int `json:"vehicle_id" validate:"required,gt=0"`
	Quantity  int `json:"quantity" validate:"required,gte=1,lte=100"`
}

// for vehicel service
type VehicleResponse struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Brand    string  `json:"brand"`
	ImageURL string  `json:"image_url"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}

type CartItemResponse struct {
	ID        int     `json:"id"`
	VehicleID int     `json:"vehicle_id"`
	Name      string  `json:"name"`
	Image     string  `json:"image"`
	Brand     string  `json:"brand"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	SubTotal  float64 `json:"subTotal"`
}
