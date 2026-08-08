package service

import (
	"context"
	"payment/internal/client"
	"payment/internal/db/storage"
	"payment/internal/models"
)

type StripeClient interface {
}

type PaymentService struct {
	Repo   storage.Storage
	Stripe StripeClient
	Order  *client.OrderClient
}

func NewPaymentService(repo storage.Storage, stripe StripeClient, order *client.OrderClient) *PaymentService {
	return &PaymentService{
		Repo:   repo,
		Stripe: stripe,
		Order:  order,
	}
}

func (p *PaymentService) CreatePayment(ctx context.Context, payment *models.Payment) error {
	return nil
}
