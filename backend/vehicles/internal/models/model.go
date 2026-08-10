package models

import "time"

type Vehicle struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name" validate:"required,min=2,max=50"`
	Model       string    `json:"model" validate:"required,min=2,max=50"`
	Price       float64   `json:"price" validate:"required"`
	Brand       string    `json:"brand" validate:"required"`
	Stock       int64     `json:"stock" validate:"required"`
	Description string    `json:"description" validate:"required"`
	ImageURL    string    `json:"image_url" validte:"required"`
	Type        string    `json:"type" validate:"required,oneof=Car Bike Truck SUV Van Bus Other"`
	Category    string    `json:"category" validate:"required,oneof=Normal Moderate Premium"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaginationResponse struct {
	Page        int64     `json:"page"`
	Limit       int64     `json:"limit"`
	Total       int64     `json:"total"`
	TotalPages  int64     `json:"total_pages"`
	HasNext     bool      `json:"has_next"`
	HasPrevious bool      `json:"has_previous"`
	Data        []Vehicle `json:"data"`
}

type VehicleQuery struct {
	Page     int
	Limit    int
	Search   string
	Type     string
	Category string
	SortBy   string
	Order    string
}

type VehicleBulkRequest struct {
	VehicleIDs []int `json:"vehicle_ids" validate:"required,min=1,dive,gt=0"`
}

type DecreaseStockRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}
