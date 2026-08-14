package db

import (
	"database/sql"
	"fmt"
	"notification/internal/config"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	data, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database : %w", err)
	}

	if err = data.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database : %w", err)
	}

	_, err = data.Exec(`
		CREATE TABLE IF NOT EXISTS notifications (
	    id BIGSERIAL PRIMARY KEY,
	    order_id BIGINT NOT NULL,
	    user_id BIGINT NOT NULL,
	    email VARCHAR(255) NOT NULL,
	    type VARCHAR(50) NOT NULL,
	    status VARCHAR(20) NOT NULL,
	    provider VARCHAR(50),
	    created_at TIMESTAMP DEFAULT NOW(),
	    sent_at TIMESTAMP
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to create table : %w", err)
	}

	return &Database{Db: data}, nil
}
