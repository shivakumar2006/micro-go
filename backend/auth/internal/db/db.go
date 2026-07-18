package db

import (
	"auth/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(cfg config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName)
	data, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db : %v", err)
	}

	if err := data.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db : %v", err)
	}

	_, err = data.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id BIGSERIAL PRIMARY KEY,
			role_name VARCHAR(20) NOT NULL UNIQUE
		)

		CREATE TABLE users (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			role_id BIGINT REFERENCES roles(id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE user_sessions (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			refresh_token_hash VARCHAR(255) NOT NULL,
			device_name VARCHAR(255),
			device_type VARCHAR(50),
			user_agent TEXT,
			ip_address VARCHAR(45),
			expires_at TIMESTAMP NOT NULL,
			revoked BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to create tables : %w", err)
	}

	return &Database{DB: data}, nil

}
