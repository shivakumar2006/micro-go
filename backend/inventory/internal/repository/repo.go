package repository

import (
	"context"
	"database/sql"
	"errors"
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
	query := `INSERT INTO inventory(
		order_id, vehicle_id, quantity, operation, status, created_at, updated_at
	) VALUES($1, $2, $3, $4, $5, NOW(), NOW())`

	_, err := r.Db.ExecContext(ctx, query, tx.OrderID, tx.VehicleID, tx.Quantity, tx.Operation, tx.Status)

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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found : %w", sql.ErrNoRows)
		}
		return nil, fmt.Errorf("failed to get transaction by id : %w", err)
	}

	return &transaction, nil
}

func (r *InventoryRepository) GetTransactionByOrderID(ctx context.Context, orderID int64) ([]models.Inventory, error) {
	var transaction []models.Inventory

	query := `SELECT id, order_id, vehicle_id, quantity, operation, status, created_at, updated_at
		FROM inventory
		WHERE order_id = $1
	`

	rows, err := r.Db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by order id : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Inventory
		err := rows.Scan(
			&t.ID,
			&t.OrderID,
			&t.VehicleID,
			&t.Quantity,
			&t.Operation,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transaction = append(transaction, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transaction, nil
}

func (r *InventoryRepository) UpdateTransactionStatus(ctx context.Context, transactionID int64, status string) error {
	query := `UPDATE inventory
	SET status = $2, updated_at = NOW()
	WHERE id = $1`

	result, err := r.Db.ExecContext(ctx, query, transactionID, status)

	if err != nil {
		return fmt.Errorf("failed to update transaction status : %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected : %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no transaction found with id : %d", transactionID)
	}
	return nil
}

func (r *InventoryRepository) GetTransactions(ctx context.Context) ([]models.Inventory, error) {
	var transactions []models.Inventory

	query := `SELECT id, order_id, vehicle_id, quantity, operation, status, created_at, updated_at
		FROM inventory
	`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions : %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var transaction models.Inventory
		if err := rows.Scan(
			&transaction.ID,
			&transaction.OrderID,
			&transaction.VehicleID,
			&transaction.Quantity,
			&transaction.Operation,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction : %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate transactions : %w", err)
	}

	return transactions, nil
}
