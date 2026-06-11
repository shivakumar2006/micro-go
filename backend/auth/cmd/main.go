package main

import (
	"auth/internal/config"
	"auth/internal/db"
	"auth/internal/pkg"
	"auth/internal/routes"
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
	// load config
	config := config.LoadConfig()

	// setup db
	storage, err := db.NewDatabase(*config)
	if err != nil {
		log.Fatal("Failed to start db : ", err)
	}
	_ = storage

	slog.Info("db started successfully", slog.String("env", config.Env), slog.String("version", "1.0.0"))

	// add routes
	router := chi.NewRouter()

	jwtManager := pkg.NewJWTManager(
		config.JWT.AccessSecret,
		config.JWT.RefreshSecret,
		config.JWT.AccessExpiry,
		config.JWT.RefreshExpiry,
	)

	authService := routes.NewAuthservice(
		storage,
		jwtManager,
		config.JWT.RefreshExpiry,
	)

	router.Post("/api/v1/auth/register", authService.Register())
	router.Post("/api/v1/auth/login", authService.Login())
	router.Post("/api/v1/auth/refresh", authService.Refresh())

	// start server

	server := &http.Server{
		Addr:    config.Server.Addr,
		Handler: router,
	}

	slog.Info("server started successfully", slog.String("address", config.Server.Addr))

	// gracefully shutdown

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", err)
		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server successfully stopped")
}
