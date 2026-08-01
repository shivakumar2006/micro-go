package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"orders/internal/config"
	"orders/internal/db"
	"orders/internal/handler"
	"orders/internal/repository"
	"orders/internal/services"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	//load cfg
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("failed to load config")
	}

	slog.Info("config loaded successfully", slog.String("env", cfg.Server.Addr))

	//db
	database, err := db.NewDatabase(*cfg)
	if err != nil {
		log.Fatalf("failed to connect to db : %v", err)
	}

	slog.Info("database connected successfully", slog.String("host", cfg.DB.Host))

	// layers
	repo, err := repository.NewRepository(database.DB)
	if err != nil {
		log.Fatalf("failed to create repo : %v", err)
	}

	service := services.NewOrderService(*repo)
	handler := handler.NewOrderHandler(service)

	// routes
	router := chi.NewRouter()

	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Timeout(10 & time.Second))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "service": "order_service"}`))
	})

	router.Group(func(r chi.Router) {
		r.Post("/api/v1/orders", handler.CreateOrder)
		r.Get("/api/v1/orders/{id}", handler.GetOrderByID)
		r.Get("/api/v1/orders/{userId}", handler.GetOrdersByUserID)
		r.Patch("/api/v1/orders/{id}/status", handler.UpdateOrderStatus)
		r.Patch("/api/v1/orders/{id}/cancel", handler.CancelOrder)
		r.Patch("/api/v1/orders/{id}/pay", handler.MarkOrderPaid)
	})

	//server

	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("server running on", slog.String("port", cfg.Server.Addr))

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		slog.Info("HTTP server start listening...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server : %v", err)
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server not stopped", slog.String("error : ", err.Error()))
		os.Exit(1)
	}

	slog.Info("server gracefully stopped")
}
