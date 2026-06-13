package main

import (
	"auth/internal/config"
	"auth/internal/db"
	"auth/internal/middleware"
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
	"github.com/go-chi/cors"
)

func main() {
	// load config
	config := config.LoadConfig()

	// setup db
	storage, err := db.NewDatabase(*config)
	if err != nil {
		log.Fatal("Failed to start db : ", err)
	}

	slog.Info("db started successfully", slog.String("env", config.Env), slog.String("version", "1.0.0"))

	// add routes
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		ExposedHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))

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

	authMiddleware := middleware.NewAuthMiddleware(*jwtManager)

	router.Post("/api/v1/auth/register", authService.Register())
	router.Post("/api/v1/auth/login", authService.Login())
	router.Post("/api/v1/auth/refresh", authService.Refresh())
	router.With(authMiddleware.Authenticate).
		Post("/api/v1/auth/logout", authService.Logout())
	router.With(authMiddleware.Authenticate).
		Post("/api/v1/auth/logoutall", authService.LogoutAll())
	router.With(authMiddleware.Authenticate).
		Get("/api/v1/auth/me", authService.GetMe())

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
