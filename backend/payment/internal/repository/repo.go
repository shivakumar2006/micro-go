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
	err := p.Db.QueryRowContext(ctx, `
		INSERT INTO payments(order_id, user_id, amount, currency, provider, payment_intent_id, stripe_session_id, status, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id
	`, payment.OrderID, payment.UserID, payment.Amount, payment.Currency, payment.Provider, payment.PaymentIntentID, payment.StripeSessionID, payment.Status).Scan(&payment.ID)

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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment by id : %v", err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment by order id : %v", err)
	}
	return &payment, nil
}

func (p *PaymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID int, status string, paymentIntentID int) error {
	result, err := p.Db.ExecContext(ctx, `
		UPDATE payments
		SET status = $2, payment_intent_id = $3, updated_at = NOW()
		WHERE id = $1
	`, paymentID, status, paymentIntentID)

	if err != nil {
		return fmt.Errorf("failed to update payment status : %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected : %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("failed to update payment status : %v", sql.ErrNoRows)
	}

	return nil
}
