package storage

import (
	"context"
	"inventory/internal/models"
)

type InventoryRepository interface {
	CreateTransaction(ctx context.Context, tx *models.Inventory) error
	GetTransactionByID(ctx context.Context, id int64) (*models.Inventory, error)
	GetTransactionByOrderID(ctx context.Context, orderID int64) (*models.Inventory, error)
	UpdateTransactionStatus(ctx context.Context, transactionID int64, status string) error
	GetTransactions(ctx context.Context) ([]models.Inventory, error)
}
