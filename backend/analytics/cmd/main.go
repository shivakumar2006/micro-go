package main

import (
	"analytics/internal/config"
	"analytics/internal/db"
	"analytics/internal/handler"
	"analytics/internal/kafka"
	"analytics/internal/repository"
	"analytics/internal/service"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	//load config
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatalf("config not found")
	}
	slog.Info("config loaded successfully", slog.Any("config", cfg))

	//db
	db, err := db.NewDatabase(*cfg)
	if err != nil {
		log.Fatalf("error while creating database connection: %v", err)
	}
	slog.Info("database connection successful", slog.String("db", cfg.DB.DBName))

	// layers
	repo := repository.NewAnalyticRepository(db.Db)

	srv := service.NewService(repo)

	handler := handler.NewHandler(srv)

	// routes
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Get("/api/v1/analytics", handler.GetPaymentAnalytic)
		r.Get("/api/v1/analytics/order/{orderId}", handler.GetPaymentByOrderID)
		r.Get("/api/v1/analytics/user/{userId}", handler.GetPaymentByUserID)
		r.Get("/api/v1/analytics/payment/{paymentId}", handler.GetPaymentByPaymentID)
	})

	// server
	server := http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("Starting analytic service", slog.String("addr", cfg.Server.Addr))

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()

	go func() {
		slog.Info("kafka consumer started")

		backoff := 1 * time.Second
		maxBackoff := 10 * time.Second

		for {
			if consumerCtx.Err() != nil {
				return
			}

			consumer := kafka.NewConsumer(
				cfg.Kafka.Addr,
				cfg.Kafka.Topic,
				cfg.Kafka.GroupID,
			)

			err := consumer.Start(consumerCtx, func(event kafka.PaymentSuccessEvent) error {
				slog.Info("payment success event received", slog.Any("event", event))

				return srv.ProcessPaymentSuccess(consumerCtx, event)
			})

			consumer.Close()

			if consumerCtx.Err() != nil {
				return
			}

			slog.Error("kafka consumer stopped", slog.String("error", err.Error()))

			timer := time.NewTimer(backoff)

			select {
			case <-timer.C:
			case <-consumerCtx.Done():
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
		slog.Info("analytics http server listening", slog.String("address", cfg.Server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start analytics service")
			log.Fatal("failed to start analytics service", "error", err)
		}
	}()

	<-quit
	consumerCancel()

	slog.Info("inventory service is shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		log.Fatal("server forced to shutdown", "error", err)
	}

	slog.Info("analytics service stopped gracefully")
}
