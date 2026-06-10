package main

import (
	"auth/internal/config"
	"auth/internal/db"
	"auth/internal/routes"
	"log"
	"log/slog"

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

	router.Post("/api/v1/auth/register", routes.Register(storage))
	// start server
}
