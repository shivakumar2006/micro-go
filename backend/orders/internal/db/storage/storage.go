package storage

import (
	"context"
	"database/sql"
	"orders/internal/models"
)

type Repository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)

	CreateOrder(ctx context.Context, tx *sql.Tx, order *models.Order) error

	CreateOrderItems(ctx context.Context, tx *sql.Tx, items []models.OrderItem) error

	GetOrderByID(ctx context.Context, id int64) (*models.Order, error)

	GetOrdersByUserID(ctx context.Context, userID int64) ([]models.Order, error)

	UpdateOrderStatus(ctx context.Context, id int64, status string) error
}
