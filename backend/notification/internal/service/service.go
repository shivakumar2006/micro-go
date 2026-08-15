package service

import (
	"context"
	"fmt"
	"notification/internal/client"
	"notification/internal/kafka"
)

type NotificationService struct {
	emailClient *client.EmailClient
}

func NewNotificationService(emailClient *client.EmailClient) *NotificationService {
	return &NotificationService{emailClient: emailClient}
}

func (s *NotificationService) HandlePaymentSuccess(ctx context.Context, event kafka.PaymentSuccessEvent) error {

	if event.OrderID <= 0 {
		return fmt.Errorf("invalid order id")
	}

	if event.UserID <= 0 {
		return fmt.Errorf("invalid user id")
	}

	if event.Email == "" {
		return fmt.Errorf("email is required")
	}

	if err := s.emailClient.SendEmail(ctx, event.Email, event.OrderID); err != nil {
		return fmt.Errorf("failed to send payment success notification: %w", err)
	}

	return nil
}
