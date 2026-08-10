package main

import (
	"context"
	"inventory/internal/client"
	"inventory/internal/config"
	"inventory/internal/db"
	"inventory/internal/handler"
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
	orderClient := client.NewOrderClient(cfg.Services.Orders.URL, retry, cb)
	vehicleClient := client.NewVehicleClient(cfg.Services.Vehicles.URL, retry, cb)

	service := service.NewInventoryService(*repo, orderClient, vehicleClient)

	handler := handler.NewInventoryHandler(service)

	// routes
	router := chi.NewRouter()

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

	go func() {
		slog.Info("HTTP server listening...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start inventory service")
			log.Fatal("failed to start inventory service", "error", err)
		}
	}()

	<-quit

	slog.Info("inventory service is shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("inventory service forced to shutdown", "error", err)
		log.Fatal("inventory service forced to shutdown", "error", err)
	}

	slog.Info("inventory service has been gracefully shut down")
}
