package repository

import (
	"context"
	"database/sql"
	"fmt"
	"inventory/internal/models"
)

type InventoryRepository struct {
	Db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{
		Db: db,
	}
}

func (r *InventoryRepository) CreateTransaction(ctx context.Context, tx *models.Inventory) error {
	if tx.Quantity <= 0 {
		return fmt.Errorf("invalid quantity : %d", tx.Quantity)
	}
	if tx.Operation != "CREDIT" && tx.Operation != "DEBIT" {
		return fmt.Errorf("invalid operation : %s", tx.Operation)
	}

	query := `INSERT INTO inventory(
		order_id, vehicle_id, quantity, operation, status, created_at, updated_at
	) VALUES($1, $2, $3, $4, $5, NOW(), NOW())`

	_, err := r.Db.ExecContext(ctx, query, tx.OrderID, tx.VehicleID, tx.Quantity, tx.Operation, tx.Status, tx.CreatedAt, tx.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create transactions : %w", err)
	}

	return nil
}

func (r *InventoryRepository) GetTransactionByID(ctx context.Context, id int64) (*models.Inventory, error) {
	var transaction models.Inventory

	query := `SELECT id, order_id, vehicle_id, quantity, operation, status, created_at, updated_at
		FROM inventory
		WHERE id = $1
	`

	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.OrderID,
		&transaction.VehicleID,
		&transaction.Quantity,
		&transaction.Operation,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by id : %w", err)
	}

	return &transaction, nil
}

func (r *InventoryRepository) GetTransactionByOrderID(ctx context.Context, orderID int64) (*models.Inventory, error) {
	var transaction models.Inventory

	query := `SELECT id, order_id, vehicle_id, quantity, operation, status, created_at, updated_at
		FROM inventory
		WHERE order_id = $1
	`

	err := r.Db.QueryRowContext(ctx, query, orderID).Scan(
		&transaction.ID,
		&transaction.OrderID,
		&transaction.VehicleID,
		&transaction.Quantity,
		&transaction.Operation,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by order id : %w", err)
	}

	return &transaction, nil
}
