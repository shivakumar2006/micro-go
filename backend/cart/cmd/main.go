package cmd

import (
	"cart/internal/config"
	"cart/internal/db"
	"log"
	"log/slog"
)

func main() {
	//load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config file : %v", err.Error())
	}

	// setup db
	database, err := db.NewCartDatabase(*cfg)
	if err != nil {
		log.Fatalf("failed to setup db : %v", err.Error())
	}
	defer database.DB.Close()

	slog.Info("Database connected successfully", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// add routes

	// server
}
