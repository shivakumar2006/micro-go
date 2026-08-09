package service

import (
	"context"
	"encoding/json"
	"fmt"
	"payment/internal/client"
	"payment/internal/db/storage"
	"payment/internal/models"

	"github.com/stripe/stripe-go/v78/webhook"
	"github.com/stripe/stripe-go/v82"
)

type PaymentService struct {
	Repo          storage.Storage
	Stripe        client.StripeClient
	Order         client.OrderClient
	WebhookSecret string
}

func NewPaymentService(repo storage.Storage, stripe client.StripeClient, order client.OrderClient, wh string) *PaymentService {
	return &PaymentService{
		Repo:          repo,
		Stripe:        stripe,
		Order:         order,
		WebhookSecret: wh,
	}
}

func (p *PaymentService) CreatePayment(ctx context.Context, req *models.CreatePaymentRequest) (*models.StripeCheckoutResponse, error) {
	if req.OrderID <= 0 {
		return nil, fmt.Errorf("invalid order id")
	}

	order, err := p.Order.GetOrder(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order : %w", err)
	}

	if order.Status != models.StatusPending {
		return nil, fmt.Errorf("order is not in pending state")
	}

	exists, err := p.Repo.ExistByOrderID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find existing order id : %w", err)
	}

	if exists {
		return nil, fmt.Errorf("payment is already created for this order")
	}

	// stripe session checkout
	checkoutSession, err := p.Stripe.CreateCheckoutSession(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session : %w", err)
	}

	payment := &models.Payment{
		OrderID:         order.ID,
		UserID:          order.UserID,
		Amount:          order.TotalAmount,
		Currency:        "INR",
		Provider:        "STRIPE",
		StripeSessionID: checkoutSession.ID,
		Status:          models.StatusPending,
	}

	if err := p.Repo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment : %w", err)
	}

	return checkoutSession, nil
}

func (p *PaymentService) GetPaymentByID(ctx context.Context, id int) (*models.Payment, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid payment id")
	}

	payment, err := p.Repo.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by id : %w", err)
	}

	return payment, nil
}

func (p *PaymentService) GetPaymentByOrderID(ctx context.Context, orderId int) (*models.Payment, error) {
	if orderId <= 0 {
		return nil, fmt.Errorf("invalid order id")
	}

	payment, err := p.Repo.GetPaymentByOrderID(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by order id : %w", err)
	}

	return payment, nil
}

func (p *PaymentService) UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error {
	if paymentID <= 0 {
		return fmt.Errorf("invalid payment id")
	}

	payment, err := p.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment by id : %w", err)
	}

	if status != models.StatusPaid {
		return fmt.Errorf("status is not in paid state")
	}

	if payment.Status == status {
		return nil
	}

	if err := p.Repo.UpdatePaymentStatus(ctx, paymentID, status); err != nil {
		return fmt.Errorf("failed to update payment status : %w", err)
	}

	if err := p.Order.UpdateOrderStatus(payment.OrderID, models.StatusPaid); err != nil {
		return fmt.Errorf("failed to update order status : %w", err)
	}

	return nil
}

func (p *PaymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	// verify signature
	if signature == "" {
		return fmt.Errorf("invalid signature")
	}

	event, err := webhook.ConstructEvent(payload, signature, p.WebhookSecret)
	if err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			return fmt.Errorf("failed to unmarshal webhook : %w", err)
		}

		payment, err := p.Repo.GetPaymentByStripeSessionID(ctx, session.ID)
		if err != nil {
			return fmt.Errorf("failed to get payment : %w", err)
		}

		if err := p.Repo.UpdatePaymentStatus(ctx, payment.ID, models.StatusPaid); err != nil {
			return fmt.Errorf("failed to update payment status : %w", err)
		}

		if err := p.Order.UpdateOrderStatus(payment.OrderID, models.StatusPaid); err != nil {
			return fmt.Errorf("failed to update order status : %w", err)
		}
	}
	return nil
}
