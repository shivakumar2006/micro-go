package main

import (
	"log"
	"log/slog"
	"payment/internal/client"
	"payment/internal/config"
	"payment/internal/db"
	"payment/internal/handler"
	"payment/internal/repository"
	"payment/internal/resilience"
	"payment/internal/service"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	// load config
	cfg := config.LoadConfig()
	if cfg == nil {
		slog.Error("failed to load config")
		log.Fatal("failed to open config")
	}
	slog.Info("Config file successfully loaded")

	// db
	database, err := db.NewDatabase(cfg)
	if err != nil {
		slog.Error("failed to initialize database")
		log.Fatalf("failed to initialize database : %v", err)
	}
	slog.Info("Database initialized successfully", slog.String("db", cfg.DB.Host))

	// layers
	repo := repository.NewPaymentRepository(database.Db)

	// retry pattern
	retry := resilience.NewRetry(3, 500*time.Millisecond, 5*time.Second, resilience.IsRetryable) // max attempts, base delay, max delay, isretryable function
	// cb pattern
	cb := resilience.NewCircuitBreaker()

	// clients
	orderClient := client.NewOrderClient(cfg.Orders.URL, retry, cb)
	stripeClient := client.NewStripeClient(cfg.Stripe.BaseURL, cfg.Stripe.SecretKey, cfg.Stripe.SuccessURL, cfg.Stripe.CancelURL, retry, cb)

	service := service.NewPaymentService(repo, *stripeClient, *orderClient)

	handler := handler.NewPaymentHandler(service)

	// routes
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Post("/api/v1/payments/create-checkout-session", handler.CreatePayment)
	})

	// server
}
