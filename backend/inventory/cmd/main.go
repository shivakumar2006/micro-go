package main

import (
	"context"
	"inventory/internal/client"
	"inventory/internal/config"
	"inventory/internal/db"
	"inventory/internal/handler"
	"inventory/internal/kafka"
	"inventory/internal/repository"
	"inventory/internal/resilience"
	"inventory/internal/service"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// load config
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("cannot load config")
	}
	slog.Info("config loaded successfully", slog.String("env", cfg.Server.Addr))

	// db
	database, err := db.NewDatabase(*cfg)
	if err != nil {
		slog.Error("failed to initialize database")
		log.Fatalf("failed to initialize database : %v", err)
	}
	slog.Info("Database connected successfully", slog.String("db", cfg.DB.Host))

	// layers
	repo := repository.NewInventoryRepository(database.DB)

	// resilience patterns
	retry := resilience.NewRetry(3, 500*time.Millisecond, 5*time.Second, resilience.IsRetryable)
	cb := resilience.NewCircuitBreaker()

	// clients
	orderClient := client.NewOrderClient(cfg.Services.Orders.URL, retry, cb, cfg.InternalServiceKey)
	vehicleClient := client.NewVehicleClient(cfg.Services.Vehicles.URL, retry, cb)

	service := service.NewInventoryService(*repo, orderClient, vehicleClient)

	handler := handler.NewInventoryHandler(service)

	// routes
	router := chi.NewRouter()

	router.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	router.Group(func(r chi.Router) {
		r.Post("/api/v1/inventory", handler.CreateTransaction)
		r.Get("/api/v1/inventory/{id}", handler.GetTransactionByID)
		r.Get("/api/v1/inventory/order/{orderId}", handler.GetTransactionByOrderID)
		r.Patch("/api/v1/inventory/{id}", handler.UpdateTransactionStatus)
		r.Get("/api/v1/inventory", handler.GetAllTransactions)
	})

	// server
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("starting inventory service", slog.String("addr", cfg.Server.Addr))

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	consumerContext, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()

	go func() {
		slog.Info("kafka consumer started")

		backoff := 1 * time.Second
		maxBackoff := 10 * time.Second

		for {
			if consumerContext.Err() != nil {
				return
			}

			consumer := kafka.NewConsumer([]string{"localhost:9092"}, "payment-success", "inventory-group")

			err := consumer.Start(consumerContext, func(event kafka.PaymentSuccessEvent) error {
				slog.Info(
					"payment success event received",
					slog.Int64("order_id", event.OrderID),
					slog.Int64("payment_id", event.PaymentID),
				)

				err := service.CreateTransaction(consumerContext, event.OrderID)
				if err != nil {
					slog.Error(
						"failed to create inventory transaction",
						slog.Int64("order_id", event.OrderID),
						slog.String("error", err.Error()),
					)
				}

				return err
			})

			consumer.Close()

			if consumerContext.Err() != nil {
				return
			}

			slog.Error("kafka consumer stopped", slog.String("error", err.Error()))

			timer := time.NewTimer(backoff)
			select {
			case <-timer.C:
			case <-consumerContext.Done():
				timer.Stop()
				return
			}

			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}()

	go func() {
		slog.Info("HTTP server listening...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start inventory service")
			log.Fatal("failed to start inventory service", "error", err)
		}
	}()

	<-quit
	consumerCancel()

	slog.Info("inventory service is shutting down gracefully")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("inventory service forced to shutdown", "error", err)
		log.Fatal("inventory service forced to shutdown", "error", err)
	}

	slog.Info("inventory service has been gracefully shut down")
}
