package main

import (
	"inventory/internal/config"
	"log"
	"log/slog"
)

func main() {
	// load config
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("cannot load config")
	}
	slog.Info("config loaded successfully", slog.String("env", cfg.Server.Addr))

	// db

	// layers

	// routes

	// server
}
