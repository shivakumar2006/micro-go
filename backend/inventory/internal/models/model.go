package models

import "time"

const (
	TransactionPending = "PENDING"
	TransactionSuccess = "SUCCESS"
	TransactionFailed  = "FAILED"
)

const (
	OperationDecrease = "DECREASE"
	OperationIncrease = "INCREASE"
)

const (
	OperationCredit = "CREDIT"
	OperationDebit  = "DEBIT"
)

type Inventory struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	VehicleID int       `json:"vehicle_id"`
	Quantity  int       `json:"quantity"`
	Operation string    `json:"operation"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderResponse struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Status      string  `json:"status"`
	TotalAmount float64 `json:"total_amount"`

	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	VehicleID int `json:"vehicle_id"`
	Quantity  int `json:"quantity"`
}

type DecreaseStockRequest struct {
	VehicleID int64 `json:"vehicle_id"`
	Quantity  int   `json:"quantity"`
}
