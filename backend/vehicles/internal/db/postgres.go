package db

import (
	"database/sql"
	"fmt"
	"vehicles/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS vehicles (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			model VARCHAR(255) NOT NULL,
			price NUMERIC(10,2) NOT NULL CHECK(price >= 0),
			stock INT NOT NULL DEFAULT 0 CHECK(price >= 0),
			brand TEXT,
			description TEXT,
			image_url TEXT,
			type VARCHAR(50) NOT NULL,

			category VARCHAR(50) NOT NULL,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	_, err = db.Exec(`
    ALTER TABLE vehicles
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
`)
	if err != nil {
		return nil, fmt.Errorf("failed to add updated_at column: %w", err)
	}

	return &Database{Db: db}, nil
}
