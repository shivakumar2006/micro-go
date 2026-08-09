package service

import (
	"context"
	"fmt"
	"payment/internal/client"
	"payment/internal/db/storage"
	"payment/internal/models"
)

type PaymentService struct {
	Repo   storage.Storage
	Stripe client.StripeClient
	Order  client.OrderClient
}

func NewPaymentService(repo storage.Storage, stripe client.StripeClient, order client.OrderClient) *PaymentService {
	return &PaymentService{
		Repo:   repo,
		Stripe: stripe,
		Order:  order,
	}
}

func (p *PaymentService) CreatePayment(ctx context.Context, payment *models.Payment) error {
	order, err := p.Order.GetOrder(payment.OrderID)
	if err != nil {
		return fmt.Errorf("failed to call order client : %w", err)
	}

	if order.Status != models.StatusPending {
		return fmt.Errorf("order status is pending")
	}

	exists, err := p.Repo.ExistByOrderID(ctx, payment.OrderID)
	if err != nil {
		return fmt.Errorf("failed to check order exist : %w", err)
	}

	if exists {
		return fmt.Errorf("payment already exists for order id : %d", payment.OrderID)
	}

	checkout, err := p.Stripe.CreateCheckoutSession(order)
	if err != nil {
		return fmt.Errorf("failed to create checkout session : %w", err)
	}

	payment = &models.Payment{
		OrderID:         order.ID,
		UserID:          order.UserID,
		Amount:          order.TotalAmount,
		Currency:        "inr",
		StripeSessionID: checkout.ID,
		Status:          models.StatusPending,
	}

	if err := p.Repo.CreatePayment(ctx, payment); err != nil {
		return fmt.Errorf("failed to create payment : %w", err)
	}

	return nil
}
