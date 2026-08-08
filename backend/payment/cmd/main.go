package cmd

import (
	"log"
	"log/slog"
	"payment/internal/config"
	"payment/internal/db"
	"payment/internal/repository"

	"github.com/go-chi/chi/v5"
)

func main() {
	// load config
	cfg := config.LoadConfig()
	if cfg == nil {
		slog.Error("failed to load config")
		log.Fatal("failed to open config")
	}
	slog.Info("Config file successfully loaded")

	// db
	database, err := db.NewDatabase(cfg)
	if err != nil {
		slog.Error("failed to initialize database")
		log.Fatalf("failed to initialize database : %v", err)
	}
	slog.Info("Database initialized successfully", slog.String("db", cfg.DB.Host))

	// layers
	repo := repository.NewPaymentRepository(database.Db)

	// routes
	router := chi.NewRouter()

	// server
}
