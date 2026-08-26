package storage

import (
	"analytics/internal/models"
	"context"
)

type Storage interface {
	CreatePaymentAnalytic(ctx context.Context, event *models.Analytics) error
	GetPaymentAnalytic(ctx context.Context) ([]models.Analytics, error)
	GetPaymentByPaymentID(ctx context.Context, id int64) (*models.Analytics, error)
	GetPaymentAnalyticsByUserID(ctx context.Context, userID int64) ([]models.Analytics, error)
	GetPaymentAnalyticsByOrderID(ctx context.Context, orderID int64) ([]models.Analytics, error)
}
