package db

import (
	"database/sql"
	"fmt"
	"orders/internal/config"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(cfg config.Config) (*Database, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db : %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			status VARCHAR(50) NOT NULL,
			total_amount NUMERIC(12,2) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to create orders table : %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id BIGSERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL,
			vehicle_id BIGINT NOT NULL,
			price NUMERIC(12, 2) NOT NULL,
			quantity INT NOT NULL,

			CONSTRAINT fk_order
			FOREIGN KEY (order_id)
			REFERENCES orders(id)
			ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to create order_items table : %v", err)
	}

	return &Database{DB: db}, nil
}
