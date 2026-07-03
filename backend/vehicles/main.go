package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vehicles/internal/config"
	"vehicles/internal/db"
	"vehicles/internal/handler"
	"vehicles/internal/middleware"
	"vehicles/internal/redis"
	"vehicles/internal/repository"
	"vehicles/internal/service"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	// setup db
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %s", err.Error())
	}
	defer database.Db.Close()

	slog.Info("database successfully initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	repo := repository.NewVehicleRepo(database.Db)
	redisCient, err := redis.NewRedis(cfg)
	if err != nil {
		log.Fatalf("failed to initialized redis : %v", err)
	}

	cache := redis.NewCache(redisCient.Client)
	service := service.NewService(repo, cache)
	handler := handler.NewVehicleHandler(service)

	// add routes
	router := chi.NewRouter()

	// global middlewares
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(10 * time.Second))

	auth := middleware.NewAuthMiddleware(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)

	router.Group(func(r chi.Router) {
		r.Use(auth.Authenticate)
		r.With(auth.RequireRole("admin")).Post("/api/v1/vehicles", handler.CreateVehicle)
		r.With(auth.RequireRole("admin")).Put("/api/v1/vehicles/{id}", handler.UpdateVehicle)
		r.With(auth.RequireRole("admin")).Delete("/api/v1/vehicles/{id}", handler.DeleteVehicle)
	})
	router.Get("/api/v1/vehicles", handler.GetAllVehicles)
	router.Get("/api/v1/vehicles/{id}", handler.GetVehicleById)

	// start server
	server := http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("server started successfully", slog.String("address", cfg.Server.Address))

	// graceful shutdown

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("HTTP server listening...")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", err)
		}
	}()

	<-quit

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Info("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server gracefully stopped")
}
