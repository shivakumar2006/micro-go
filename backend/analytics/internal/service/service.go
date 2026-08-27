package service

import (
	"analytics/internal/db/storage"
	"analytics/internal/kafka"
	"analytics/internal/models"
	"context"
	"fmt"
)

type Service struct {
	storage storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) ProcessPaymentSuccess(ctx context.Context, event kafka.PaymentSuccessEvent) error {
	analytics := models.Analytics{
		PaymentID: event.PaymentID,
		OrderID:   event.OrderID,
		UserID:    event.UserID,
		Email:     event.Email,
		Status:    event.Status,
		PaidAt:    event.CreatedAt,
		CreatedAt: event.CreatedAt,
	}

	if err := s.storage.CreatePaymentAnalytic(ctx, &analytics); err != nil {
		return fmt.Errorf("failed to insert payment analytic : %w", err)
	}

	return nil
}

func (s *Service) GetPaymentAnalytics(ctx context.Context) ([]models.Analytics, error) {
	data, err := s.storage.GetPaymentAnalytics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment analytics : %w", err)
	}
	return data, nil
}

func (s *Service) GetPaymentByPaymentID(ctx context.Context, paymentID int64) (*models.Analytics, error) {
	data, err := s.storage.GetPaymentByPaymentID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by payment id : %w", err)
	}
	return data, nil
}

func (s *Service) GetPaymentByOrderID(ctx context.Context, orderID int64) (*models.Analytics, error) {
	data, err := s.storage.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by order id : %w", err)
	}
	return data, nil
}

func (s *Service) GetPaymentByUserID(ctx context.Context, userID int64) ([]models.Analytics, error) {
	data, err := s.storage.GetPaymentByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by user id : %w", err)
	}
	return data, nil
}
