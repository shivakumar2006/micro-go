package storage

import (
	"context"
	"inventory/internal/models"
)

type InventoryRepository interface {

	// Inventory transaction create karne ke liye
	CreateTransaction(ctx context.Context, tx *models.Inventory) error

	// Order already process hua ya nahi
	GetTransactionByOrderID(ctx context.Context, orderID int64) (*models.Inventory, error)

	// Status update
	UpdateTransactionStatus(ctx context.Context, transactionID int64, status string) error

	// By ID
	GetTransactionByID(ctx context.Context, id int64) (*models.Inventory, error)

	// (Optional) Admin / Debugging
	GetTransactions(ctx context.Context) ([]models.Inventory, error)
}
