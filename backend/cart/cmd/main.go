package main

import (
	"cart/internal/client"
	"cart/internal/config"
	"cart/internal/db"
	"cart/internal/handler"
	"cart/internal/metrics"
	"cart/internal/middleware"
	"cart/internal/redis"
	"cart/internal/repository"
	"cart/internal/resilience"
	"cart/internal/service"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	//load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config file : %v", err.Error())
	}

	// setup db
	database, err := db.NewCartDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to setup db : %v", err.Error())
	}
	defer database.DB.Close()

	slog.Info("Database connected successfully", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// layers
	repo := repository.NewCartRepository(database.DB)
	redisClient, err := redis.NewRedis(cfg)
	if err != nil {
		log.Fatalf("failed to initialized redis : %v", err)
	}

	slog.Info("Redis connected successfully", slog.String("env", cfg.Env))

	defer redisClient.Close()

	cache := redis.NewCache(redisClient.Client)

	// retry
	retry := resilience.NewRetry(3, 500*time.Millisecond, 5*time.Second, resilience.IsRetryable)
	// circuit breaker
	cb := resilience.NewCircuitBreaker()
	vehicleClient := client.NewVehicleClient(cfg.Vechile.URL, retry, cb)

	service := service.NewCartService(repo, cache, vehicleClient)
	handler := handler.NewCartHandler(service)

	// add routes
	router := chi.NewRouter()

	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(10 * time.Second))

	router.Use(metrics.MerticsMiddleware)

	auth := middleware.NewAuthMiddleware(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)

	router.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	router.Group(func(r chi.Router) {
		r.Use(auth.Authenticate)
		r.With(auth.RequireRole("customer")).Post("/api/v1/cart", handler.AddToCart)
		r.With(auth.RequireRole("customer")).Put("/api/v1/cart/{id}", handler.UpdateCartQuantity)
		r.With(auth.RequireRole("customer")).Get("/api/v1/cart", handler.GetUserCart)
		r.With(auth.RequireRole("customer")).Delete("/api/v1/cart/{id}", handler.DeleteCartItem)
		r.With(auth.RequireRole("customer")).Delete("/api/v1/cart", handler.DeleteCart)
		r.With(auth.RequireRole("customer")).Get("/api/v1/cart/total", handler.GetCartTotal)
		r.With(auth.RequireRole("customer")).Get("/api/v1/cart/count", handler.CountItems)
	})

	// server
	server := http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("Server started successfully", slog.String("address : ", cfg.Server.Address))

	quit := make(chan (os.Signal), 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		slog.Info("HTTP server listening.....")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", err)
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("server shutting down gracefully...")

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server not stopped ", slog.String("error : ", err.Error()))
		os.Exit(1)
	}

	slog.Info("server gracefully stopped")
}
