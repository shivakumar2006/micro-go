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
