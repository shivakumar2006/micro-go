package main

import (
	"inventory/internal/client"
	"inventory/internal/config"
	"inventory/internal/db"
	"inventory/internal/handler"
	"inventory/internal/repository"
	"inventory/internal/resilience"
	"inventory/internal/service"
	"log"
	"log/slog"
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
		slog.Error("failed to initialize database : %v", err)
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
	})

	// server
}
