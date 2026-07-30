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
	err := tx.QueryRowContext(ctx, `
		INSERT INTO orders(user_id, total_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, order.UserID, order.TotalAmount, order.Status, order.CreatedAt, order.UpdatedAt).Scan(&order.ID)

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

func (r *Repository) CreateOrderItem(ctx context.Context, tx *sql.Tx, items []models.OrderItem) error {
	query := `
		INSERT INTO order_items(order_id, vehicle_id, quantity, price, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6)
	`

	for _, item := range items {
		_, err := tx.ExecContext(ctx, query, item.OrderID, item.VehicleID, item.Quantity, item.Price, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23503" {
					return fmt.Errorf("invalid order id : %w", err)
				}
			}
			return fmt.Errorf("failed to create order_item : %w", err)
		}
	}

	return nil
}

func (r *Repository) GetOrderByID(ctx context.Context, id int64) (*models.Order, error) {
	var order models.Order

	err := r.DB.QueryRowContext(ctx, `
		SELECT id, user_id, total_amount, status, created_at, updated_at
		FROM orders 
		WHERE id = $1
	`, id).Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order not found : %w", err)
		}
		return nil, fmt.Errorf("failed to get order : %w", err)
	}

	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, order_id, vehicle_id, quantity, price, created_at, updated_at
		FROM order_items
		WHERE order_id = $1
	`, order.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get order items : %w", err)
	}

	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem

		err := rows.Scan(&item.ID, &item.OrderID, &item.VehicleID, &item.Quantity, &item.Price, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item : %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get order items : %w", err)
	}

	order.Items = items

	return &order, nil
}

func (r *Repository) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	var orders []models.Order

	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, user_id, total_amount, status, created_at, updated_at 
		FROM orders 
		WHERE user_id = $1
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get orders : %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var order models.Order

		err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order : %w", err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return orders, nil
	}

	var orderIDs []int
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	itemsRow, err := r.DB.QueryContext(ctx, `
		SELECT id, order_id, vehicle_id, quantity, price, created_at, updated_at
		FROM order_items
		WHERE order_id = ANY($1)
	`, orderIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to get order items : %w", err)
	}

	defer itemsRow.Close()

	for itemsRow.Next() {
		var items models.OrderItem

		err := itemsRow.Scan(&items.ID, &items.OrderID, &items.VehicleID, &items.Quantity, &items.Price, &items.CreatedAt, &items.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order items : %w", err)
		}

		for i := range orders {
			if orders[i].ID == items.OrderID {
				orders[i].Items = append(orders[i].Items, items)
				break
			}
		}
	}

	if err := itemsRow.Err(); err != nil {
		return nil, fmt.Errorf("failed to get order items : %w", err)
	}

	return orders, nil
}
