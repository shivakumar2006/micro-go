package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"orders/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.DB.BeginTx(ctx, nil)
}

func (r *Repository) CreateOrder(ctx context.Context, tx *sql.Tx, order *models.Order) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO orders(user_id, total_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, order.UserID, order.TotalAmount, order.Status, order.CreatedAt, order.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return fmt.Errorf("invalid user : %w", err)
			}
		}
		return fmt.Errorf("failed to create order : %w", err)
	}

	return nil
}

func (r *Repository) CreateOrderItem(ctx context.Context, tx *sql.Tx, item []models.OrderItems) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO order_items(order_id, vehicle_id, quantity, price, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6)
	`, item.OrderID, item.VehicleID, item.Quantity, item.Price, item.CreatedAt, item.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return fmt.Errorf("invalid order : %w", err)
			}
		}
		return fmt.Errorf("failed to create order_item : %w", err)
	}

	return nil
}

func (r *Repository) GetOrders(ctx context.Context)
