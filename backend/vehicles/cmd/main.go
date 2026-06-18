package main

import (
	"log"
	"vehicles/internal/config"
	"vehicles/internal/db"
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

	log.Println("Database initialized successfully")

	// add routes
	// start server
}
