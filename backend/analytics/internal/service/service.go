package service

import (
	"analytics/internal/db/storage"
	"analytics/internal/kafka"
	"context"
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

}
