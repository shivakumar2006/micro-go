package storage

import (
	"analytics/internal/models"
	"context"
)

type Storage interface {
	CreatePaymentAnalytic(ctx context.Context, event *models.Analytics) error
}
