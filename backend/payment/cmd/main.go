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
	"payment/internal/kafka"
	"payment/internal/middleware"
	"payment/internal/repository"
	"payment/internal/resilience"
	"payment/internal/service"
	"payment/internal/worker"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	outboxRepo := repository.NewOutboxEventRepository(database.Db)

	// retry pattern
	retry := resilience.NewRetry(3, 500*time.Millisecond, 5*time.Second, resilience.IsRetryable) // max attempts, base delay, max delay, isretryable function
	// cb pattern
	cb := resilience.NewCircuitBreaker()

	// clients
	orderClient := client.NewOrderClient(cfg.Orders.URL, retry, cb, cfg.InternalServiceKey)
	stripeClient := client.NewStripeClient(cfg.Stripe.BaseURL, cfg.Stripe.SecretKey, cfg.Stripe.SuccessURL, cfg.Stripe.CancelURL, retry, cb)

	// kafka
	producer := kafka.NewProducer([]string{"localhost:9092"}, "payment-success")
	defer producer.Close()

	outboxWorker := worker.NewOutboxWorker(outboxRepo, producer)

	service := service.NewPaymentService(repo, outboxRepo, *stripeClient, *orderClient, cfg.Stripe.WebhookSecret)

	handler := handler.NewPaymentHandler(service)

	// routes
	router := chi.NewRouter()

	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(10 * time.Second))

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.AccessTokenSecret, cfg.JWT.RefreshTokenSecret)

	// prometheus
	router.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	router.Post("/api/v1/payments/webhook", handler.WebhookHandler)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		r.With(authMiddleware.RequireRole("customer")).Post("/api/v1/payments/create-checkout-session", handler.CreatePayment)
		r.With(authMiddleware.RequireRole("customer")).Get("/api/v1/payments/{id}", handler.GetPaymentByID)
		r.With(authMiddleware.RequireRole("customer")).Get("/api/v1/payments/order/{orderid}", handler.GetPaymentByOrderID)
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

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		slog.Info("outbox worker started")

		for {
			select {
			case <-ticker.C:
				slog.Info("outbox worker tick")

				if err := outboxWorker.ProcessPendingEvents(workerCtx); err != nil {
					slog.Error(
						"failed to process pending event",
						slog.String("error", err.Error()),
					)
				}

				slog.Info("outbox worker cycle completed")

			case <-workerCtx.Done():
				slog.Info("outbox worker stopped")
				return
			}
		}
	}()

	go func() {
		slog.Info("HTTP server is starting...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", slog.String("error", err.Error()))
			log.Fatal("failed to start server")
		}
	}()

	<-quit
	workerCancel()

	slog.Info("payment service is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("payment service gracefully stopped")
}
