package repository

import (
	"analytics/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type AnalyticRepository struct {
	Db *sql.DB
}

func NewAnalyticRepository(db *sql.DB) *AnalyticRepository {
	return &AnalyticRepository{Db: db}
}

func (a *AnalyticRepository) CreatePaymentAnalytic(ctx context.Context, event *models.Analytics) error {
	query := `INSERT INTO payment_analytics(payment_id, order_id, user_id, email, status, paid_at, created_at)
			VALUES($1, $2, $3, $4, $5, $6, NOW())
	`
	_, err := a.Db.ExecContext(ctx, query, event.PaymentID, event.OrderID, event.UserID, event.Email, event.Status, event.PaidAt)
	if err != nil {
		return fmt.Errorf("failed to create payment analytics : %w", err)
	}

	return nil
}

func (a *AnalyticRepository) GetPaymentAnalytics(ctx context.Context) ([]models.Analytics, error) {
	var payments []models.Analytics

	query := `SELECT id, payment_id, order_id, user_id, email, status, paid_at, created_at
			FROM payment_analytics
	`

	rows, err := a.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment analytics : %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var payment models.Analytics
		if err := rows.Scan(&payment.ID, &payment.PaymentID, &payment.OrderID, &payment.UserID, &payment.Email, &payment.Status, &payment.PaidAt, &payment.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan rows : %w", err)
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate transactions : %w", err)
	}

	return payments, nil
}

func (a *AnalyticRepository) GetPaymentByPaymentID(ctx context.Context, paymentID int64) (*models.Analytics, error) {
	query := `SELECT id, payment_id, order_id, user_id, email, status, paid_at, created_at
			FROM payment_analytics 
			WHERE payment_id = $1
	`
	var payment models.Analytics
	err := a.Db.QueryRowContext(ctx, query, paymentID).Scan(&payment.ID, &payment.PaymentID, &payment.OrderID, &payment.UserID, &payment.Email, &payment.Status, &payment.PaidAt, &payment.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("rows not found : %w", err)
		}
		return nil, fmt.Errorf("failed to find payment analytic : %w", err)
	}

	return &payment, nil
}

func (a *AnalyticRepository) GetPaymentByUserID(ctx context.Context, userID int64) ([]models.Analytics, error) {
	query := `SELECT id, payment_id, order_id, user_id, email, status, paid_at, created_at
			FROM payment_analytics
			WHERE user_id = $1
			ORDER BY paid_at DESC
	`

	var payments []models.Analytics
	rows, err := a.Db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find payment analytics by user id : %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var payment models.Analytics
		if err := rows.Scan(&payment.ID, &payment.PaymentID, &payment.OrderID, &payment.UserID, &payment.Email, &payment.Status, &payment.PaidAt, &payment.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan rows : %w", err)
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows : %w", err)
	}

	return payments, nil
}

func (a *AnalyticRepository) GetPaymentByOrderID(ctx context.Context, orderID int64) (*models.Analytics, error) {
	query := `SELECT id, payment_id, order_id, user_id, email, status, paid_at, created_at
			FROM payment_analytics
			WHERE order_id = $1
	`

	var payment models.Analytics
	err := a.Db.QueryRowContext(ctx, query, orderID).Scan(&payment.ID, &payment.PaymentID, &payment.OrderID, &payment.UserID, &payment.Email, &payment.Status, &payment.PaidAt, &payment.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("rows not found : %w", err)
		}
		return nil, fmt.Errorf("failed to find payment analytics by order id : %w", err)
	}

	return &payment, nil
}
