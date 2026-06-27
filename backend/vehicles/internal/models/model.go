package models

import "time"

type Vehicle struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name" validate:"required,min=2,max=50"`
	Model     string    `json:"model" validate:"required,min=2,max=50"`
	Type      string    `json:"type" validate:"required,oneof=Car Bike Truck SUV Van Bus Other"`
	Category  string    `json:"category" validate:"required,oneof=Normal Moderate Premium"`
	CreatedAt time.Time `json:"created_at"`
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
