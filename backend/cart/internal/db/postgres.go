package db

import (
	"cart/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type CartDatabase struct {
	DB *sql.DB
}

func NewCartDatabase(cfg *config.Config) (*CartDatabase, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db : %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db : %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cart(
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			vehicle_id BIGINT NOT NULL,
			quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
			price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

			CONSTRAINT unique_user_vehicle UNIQUE (user_id, vehicle_id)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table : %w", err)
	}

	return &CartDatabase{DB: db}, nil
}
