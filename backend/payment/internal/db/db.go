package db

import (
	"database/sql"
	"fmt"
	"payment/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase(config *config.Config) (*Database, error) {
	data, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db : %w", err)
	}

	if err = data.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db : %w", err)
	}

	_, err = data.Exec(`
		CREATE TABLE IF NOT EXISTS payments(
			id BIGSERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			amount NUMERIC(12, 2) NOT NULL,
			currency VARCHAR(10) NOT NULL DEFAULT 'INR',
			provider VARCHAR(50) NOT NULL DEFAULT 'STRIPE',
			payment_intent_id VARCHAR(100) UNIQUE,
			stripe_session_id VARCHAR(100) UNIQUE,
			status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to create db : %w", err)
	}

	return &Database{Db: data}, nil
}
