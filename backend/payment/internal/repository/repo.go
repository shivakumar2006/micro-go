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

func (p *PaymentRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return p.Db.BeginTx(ctx, nil)
}

func (p *PaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	err := p.Db.QueryRowContext(ctx, `
		INSERT INTO payments(order_id, user_id, email, amount, currency, provider, payment_intent_id, stripe_session_id, status, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id
	`, payment.OrderID, payment.UserID, payment.Email, payment.Amount, payment.Currency, payment.Provider, payment.PaymentIntentID, payment.StripeSessionID, payment.Status).Scan(&payment.ID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf(
					"unique constraint violation: constraint=%s detail=%s",
					pgErr.ConstraintName,
					pgErr.Detail,
				)
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
		SELECT id, order_id, user_id, email, amount, currency, provider, payment_intent_id, stripe_session_id, status, created_at, updated_at
		FROM payments
		WHERE id = $1
	`, id).Scan(&payment.ID, &payment.OrderID, &payment.UserID, &payment.Email, &payment.Amount, &payment.Currency, &payment.Provider, &payment.PaymentIntentID, &payment.StripeSessionID, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)

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
		SELECT id, order_id, user_id, email, amount, currency, provider, payment_intent_id, stripe_session_id, status, created_at, updated_at
		FROM payments
		WHERE order_id = $1
	`, orderID).Scan(&payment.ID, &payment.OrderID, &payment.UserID, &payment.Email, &payment.Amount, &payment.Currency, &payment.Provider, &payment.PaymentIntentID, &payment.StripeSessionID, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment by order id : %v", err)
	}
	return &payment, nil
}

func (p *PaymentRepository) UpdatePaymentStatusTx(ctx context.Context, tx *sql.Tx, paymentID int, status string) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE payments
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`, paymentID, status)

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

func (p *PaymentRepository) ExistByOrderID(ctx context.Context, orderID int) (bool, error) {
	var count int
	err := p.Db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM payments WHERE order_id = $1
	`, orderID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check exist by order id : %v", err)
	}
	return count > 0, nil
}

func (p *PaymentRepository) GetPaymentByStripeSessionID(ctx context.Context, sessionID string) (*models.Payment, error) {
	var payment models.Payment

	err := p.Db.QueryRowContext(ctx, `
		SELECT id, order_id, user_id, email, amount, currency, provider, payment_intent_id, stripe_session_id, status, created_at, updated_at
		FROM payments
		WHERE stripe_session_id = $1
	`, sessionID).Scan(&payment.ID, &payment.OrderID, &payment.UserID, &payment.Email, &payment.Amount, &payment.Currency, &payment.Provider, &payment.PaymentIntentID, &payment.StripeSessionID, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, fmt.Errorf("failed to get payment by stripe session id : %v", err)
	}
	return &payment, nil
}
