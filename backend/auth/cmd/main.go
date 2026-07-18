package main

import (
	"auth/internal/config"
	"auth/internal/db"
	"auth/internal/handler"
	"auth/internal/middleware"
	"auth/internal/pkg"
	"auth/internal/repository"
	"auth/internal/services"
	"log"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
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
		log.Fatal("failed to connect to database : %v", err)
	}

	defer database.DB.Close()

	slog.Info("Database connected successfully", slog.String("env", cfg.ENV))

	// routes
	router := chi.NewRouter()

	accessExpiry, err := time.ParseDuration(cfg.JWT.AccessExpiry)
	if err != nil {
		log.Fatal("failed to parse access expiry: ", err)
	}

	refreshExpiry, err := time.ParseDuration(cfg.JWT.RefreshExpiry)
	if err != nil {
		log.Fatal("failed to parse refresh expiry: ", err)
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

	authRepo, err := repository.NewAuthRepository(database.DB)
	if err != nil {
		log.Fatal("failed to initialize AuthRepository: ", err)
	}

	authService := services.NewAuthService(authRepo, jwtManager, refreshExpiry)

	authHandler := handler.NewAuthHandler(authService)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	router.Post("/api/v1/auth/register", authHandler.Register())
	router.Post("/api/v1/auth/login", authHandler.Login())
	router.Post("/api/v1/auth/refresh", authHandler.Refresh())

	// server
}
