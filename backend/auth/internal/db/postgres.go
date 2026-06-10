package db

import (
	"auth/internal/config"
	"auth/internal/model"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase(cfg config.Config) (*Database, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize table: %w", err)
	}

	log.Println("Database connected successfully")

	return &Database{Db: db}, nil
}

func (d *Database) CreateUser(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := d.Db.QueryRowContext(ctx, `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, created_at
	`, user.Name, user.Email, user.Password).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if query != nil {
		return fmt.Errorf("failed to create user : %s", query.Error())
	}

	return nil
}
