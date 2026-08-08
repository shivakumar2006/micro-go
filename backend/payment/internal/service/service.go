package service

import "payment/internal/db/storage"

type PaymentService struct {
	Repo   storage.Storage
	Stripe StripeClient
}
