package db

import (
	"analytics/internal/config"
	"database/sql"
	"fmt"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase(config config.Config) (*Database, error) {
	data, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if err := data.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database : %w", err)
	}

	_, err = data.Exec(`
		CREATE TABLE IF NOT EXISTS payment_analytics(
			id BIGSERIAL PRIMARY KEY,
			payment_id BIGINT NOT NULL UNIQUE,
			order_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			email VARCHAR(200) NOT NULL,
			status VARCHAR(50) NOT NULL,
			paid_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_payment_analytics_user_id
			ON payment_analytics(user_id);

		CREATE INDEX IF NOT EXISTS idx_payment_analytics_paid_at
			ON payment_analytics(paid_at);

		CREATE INDEX IF NOT EXISTS idx_payment_analytics_order_id
			ON payment_analytics(order_id)
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to create db : %w", err)
	}

	return &Database{Db: data}, nil
}
