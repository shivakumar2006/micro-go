package storage

import (
	"context"
	"database/sql"
	"payment/internal/models"
)

type Storage interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByID(ctx context.Context, id int) (*models.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID int) (*models.Payment, error)
	UpdatePaymentStatusTx(ctx context.Context, tx *sql.Tx, paymentID int, status string) error
	ExistByOrderID(ctx context.Context, orderID int) (bool, error)
	GetPaymentByStripeSessionID(ctx context.Context, sessionID string) (*models.Payment, error)

	// outbox
	SaveEventTx(ctx context.Context, tx *sql.Tx, event *models.OutboxEvent) error
	GetPendingEvents(ctx context.Context, status string) ([]models.OutboxEvent, error)
	MarkAsPublished(ctx context.Context, eventID int) error
}
