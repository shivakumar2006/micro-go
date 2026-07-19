package main

import (
	"auth/internal/config"
	"auth/internal/db"
	"auth/internal/handler"
	"auth/internal/middleware"
	"auth/internal/pkg"
	"auth/internal/repository"
	"auth/internal/services"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// load config
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("failed to load config")
	}

	// db
	database, err := db.NewDatabase(*cfg)
	if err != nil {
		log.Fatalf("failed to connect to database : %v", err)
	}

	defer database.DB.Close()

	slog.Info("Database connected successfully", slog.String("env", cfg.ENV))

	// routes
	router := chi.NewRouter()

	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(10 * time.Second))

	accessExpiry, err := time.ParseDuration(cfg.JWT.AccessExpiry)
	if err != nil {
		log.Fatalf("failed to parse access expiry: %v", err)
	}

	refreshExpiry, err := time.ParseDuration(cfg.JWT.RefreshExpiry)
	if err != nil {
		log.Fatalf("failed to parse refresh expiry: %v", err)
	}

	jwtManager, err := pkg.NewJWTManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		accessExpiry,
		refreshExpiry,
	)
	if err != nil {
		log.Fatal("failed to initialize JWT manager: ", err)
	}

	repo, err := repository.NewAuthRepository(database.DB)
	if err != nil {
		log.Fatal("failed to initialize AuthRepository: ", err)
	}

	service := services.NewAuthService(repo, jwtManager, refreshExpiry)

	handler := handler.NewAuthHandler(service)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		r.Get("/api/v1/auth/me", handler.GetMe)
		r.Post("/api/v1/auth/logout", handler.Logout)
		r.Post("/api/v1/auth/logout-all", handler.LogoutAll)
	})
	router.Post("/api/v1/auth/register", handler.Register)
	router.Post("/api/v1/auth/login", handler.Login)
	router.Post("/api/v1/auth/refresh", handler.Refresh)

	// server
	server := http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("server started successfully", slog.String("address", cfg.Server.Addr))

	// gracefully shutdonw

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("HTTP server listening....")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-quit

	slog.Info("shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shitdown", slog.String("error", err.Error()))
	}

	slog.Info("Server gracefully stopped")
}
