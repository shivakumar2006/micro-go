package main

import (
	"log"
	"log/slog"
	"vehicles/internal/config"
	"vehicles/internal/db"

	"github.com/go-chi/chi/v5"
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

	// add routes
	_ = chi.NewRouter()

	// start server
}
