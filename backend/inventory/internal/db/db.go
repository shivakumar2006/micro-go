package db

import (
	"database/sql"
	"fmt"
	"inventory/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(cfg config.Config) (*Database, error) {
	data, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database : %w", err)
	}

	if err := data.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database : %w", err)
	}

	_, err = data.Exec(`
		CREATE TABLE IF NOT EXISTS inventory(
			id BIGSERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL,
			vehicle_id BIGINT NOT NULL,
			quantity INT NOT NULL,
			operation VARCHAR(20) NOT NULL,
			status VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to create inventory table : %w", err)
	}

	return &Database{DB: data}, nil
}
