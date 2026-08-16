package service

import (
	"context"
	"encoding/json"
	"fmt"
	"payment/internal/client"
	"payment/internal/db/storage"
	"payment/internal/kafka"
	"payment/internal/models"
	"time"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type PaymentService struct {
	Repo          storage.Storage
	Stripe        client.StripeClient
	Order         client.OrderClient
	WebhookSecret string
	Producer      *kafka.Producer
}

func NewPaymentService(repo storage.Storage, stripe client.StripeClient, order client.OrderClient, wh string, producer *kafka.Producer) *PaymentService {
	return &PaymentService{
		Repo:          repo,
		Stripe:        stripe,
		Order:         order,
		WebhookSecret: wh,
		Producer:      producer,
	}
}

func (p *PaymentService) CreatePayment(ctx context.Context, req *models.CreatePaymentRequest, authHeader string, userEmail string) (*models.StripeCheckoutResponse, error) {
	if req.OrderID <= 0 {
		return nil, fmt.Errorf("invalid order id")
	}

	order, err := p.Order.GetOrder(req.OrderID, authHeader)
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
		Email:           userEmail,
		Amount:          order.TotalAmount,
		Currency:        "INR",
		Provider:        "STRIPE",
		StripeSessionID: checkoutSession.ID,
		PaymentIntentID: nil,
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

	// transaction
	tx, err := p.Repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction : %w", err)
	}

	// if anything fails before commit
	defer tx.Rollback()

	if err := p.Repo.UpdatePaymentStatusTx(ctx, tx, paymentID, status); err != nil {
		return fmt.Errorf("failed to update payment status : %w", err)
	}

	// kafka event
	event := kafka.PaymentSuccessEvent{
		OrderID:   int64(payment.OrderID),
		PaymentID: int64(paymentID),
		UserID:    int64(payment.UserID),
		Email:     payment.Email,
		Status:    models.StatusPaid,
		CreatedAt: payment.CreatedAt,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal payment success event : %w", err)
	}

	outboxEvent := &models.OutboxEvent{
		AggregateType: "payment",
		AggregateID:   int64(payment.ID),
		EventType:     "payment-success",
		Payload:       payload,
		Status:        models.OutboxStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// save event transaction in outbox_event table
	if err := p.Repo.SaveEventTx(ctx, tx, outboxEvent); err != nil {
		return fmt.Errorf("failed to save event in outbox : %w", err)
	}

	// commit transac
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}

	if err := p.Order.UpdateOrderStatus(payment.OrderID, models.StatusPaid); err != nil {
		return fmt.Errorf("failed to update order status : %w", err)
	}

	if err := p.Producer.Publish(ctx, event); err != nil {
		return fmt.Errorf("failed to publish payment success event : %w", err)
	}

	return nil
}

func (p *PaymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	// verify signature
	if signature == "" {
		return fmt.Errorf("invalid signature")
	}

	event, err := webhook.ConstructEventWithOptions(
		payload,
		signature,
		p.WebhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)

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

		if err := p.UpdatePaymentStatus(ctx, payment.ID, models.StatusPaid); err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}
	}
	return nil
}
