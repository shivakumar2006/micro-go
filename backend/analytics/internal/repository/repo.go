package repository

import (
	"analytics/internal/models"
	"context"
	"database/sql"
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
	_, err := a.Db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create payment analytics table : %w", err)
	}

	return nil
}

func (a *AnalyticRepository) GetPaymentAnalytic(ctx context.Context) ([]models.Analytics, error) {
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
