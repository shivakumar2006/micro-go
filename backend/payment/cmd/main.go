package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"payment/internal/client"
	"payment/internal/config"
	"payment/internal/db"
	"payment/internal/handler"
	"payment/internal/repository"
	"payment/internal/resilience"
	"payment/internal/service"
	"syscall"
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

	service := service.NewPaymentService(repo, *stripeClient, *orderClient, cfg.Stripe.WebhookSecret)

	handler := handler.NewPaymentHandler(service)

	// routes
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Post("/api/v1/payments/create-checkout-session", handler.CreatePayment)
		r.Get("/api/v1/payments/{id}", handler.GetPaymentByID)
		r.Get("/api/v1/payments/order/{orderid}", handler.GetPaymentByOrderID)
		r.Post("/api/v1/payments/webhook", handler.WebhookHandler)
	})

	// server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s", cfg.Server.Addr),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	slog.Info("Server running on port", slog.String("address", cfg.Server.Addr))

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		slog.Info("HTTP server is starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", slog.String("error", err.Error()))
			log.Fatal("failed to start server")
		}
	}()

	<-quit

	slog.Info("payment service is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("payment service gracefully stopped")
}
