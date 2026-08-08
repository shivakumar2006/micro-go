package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"payment/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

type PaymentRepository struct {
	Db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{
		Db: db,
	}
}

func (p *PaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	_, err := p.Db.ExecContext(ctx, `
		INSERT INTO payments(order_id, user_id, amount, currency, provider, payment_intent_id, stripe_session_id, status, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`, payment.OrderID, payment.UserID, payment.Amount, payment.Currency, payment.Provider, payment.PaymentIntentID, payment.StripeSessionID, payment.Status)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("order already exists")
			}
			return fmt.Errorf("database error : %v", pgErr.Message)
		}
		return fmt.Errorf("failed to create payment : %v", err)
	}

	return nil
}

func (p *PaymentRepository) GetPaymentByID(ctx context.Context, id int) (*models.Payment, error) {
	var payment models.Payment

	err := p.Db.QueryRowContext(ctx, `
		SELECT id, order_id, user_id, amount, currency, provider, payment_intent_id, stripe_session_id, status, created_at, updated_at
		FROM payments
		WHERE id = $1
	`, id).Scan(&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount, &payment.Currency, &payment.Provider, &payment.PaymentIntentID, &payment.StripeSessionID, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("payment not found: %v", err)
	}
	return &payment, nil
}

func (p *PaymentRepository) GetPaymentByOrderID(ctx context.Context, orderID int) (*models.Payment, error) {
	var payment models.Payment

	err := p.Db.QueryRowContext(ctx, `
		SELECT id, order_id, user_id, amount, currency, provider, payment_intent_id, stripe_session_id, status, created_at, updated_at
		FROM payments
		WHERE order_id = $1
	`, orderID).Scan(&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount, &payment.Currency, &payment.Provider, &payment.PaymentIntentID, &payment.StripeSessionID, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("payment not found: %v", err)
	}
	return &payment, nil
}

func (p *PaymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error {

}
