package models

import "time"

type Vehicle struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Model     string    `json:"model"`
	Type      string    `json:"type"`
	Category  string    `json:"category"`
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
