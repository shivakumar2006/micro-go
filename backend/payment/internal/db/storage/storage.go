package storage

import (
	"context"
	"payment/internal/models"
)

type Storage interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByID(ctx context.Context, id int) (*models.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID int) (*models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error
	ExistByOrderID(ctx context.Context, orderID int) (bool, error)
}
