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
		return nil, fmt.Errorf("failed to initialize user table: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			token TEXT NOT NULL UNIQUE,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (user_id)
			REFERENCES users(id)
			ON DELETE CASCADE
		)		
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize refresh token table: %w", err)
	}

	log.Println("Database connected successfully")

	return &Database{Db: db}, nil
}

func (d *Database) CreateUser(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := d.Db.QueryRowContext(ctx, `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, created_at
	`, user.Name, user.Email, user.Password).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user : %s", err.Error())
	}

	return nil
}

func (d *Database) FindUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User

	err := d.Db.QueryRowContext(ctx, `
		SELECT id, name, email, password, created_at 
		FROM users
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("failed to find user by email : %s", err.Error())
	}

	return &user, nil
}

func (d *Database) FindUserById(id int64) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User

	err := d.Db.QueryRowContext(ctx, `
		SELECT id, name, email, password, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to find user by id : %s", err.Error())
	}

	return &user, nil
}

// REFRESH TOKEN QUERIES
func (d *Database) SaveRefreshToken(token *model.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token.CreatedAt = time.Now()

	err := d.Db.QueryRowContext(ctx, `
		INSERT INTO refresh_tokens(user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, token, expires_at, created_at
	`, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt).Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to save refresh token : %s", err.Error())
	}

	return nil
}

func (d *Database) FindRefreshToken(tokenString string) (*model.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var refreshToken model.RefreshToken
	err := d.Db.QueryRowContext(ctx, `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1
	`, tokenString).Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to find refresh token : %s", err.Error())
	}

	return &refreshToken, nil
}

func (d *Database) DeleteRefreshToken(tokenString string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := d.Db.ExecContext(ctx, `
		DELETE FROM refresh_tokens
		WHERE token = $1 
	`, tokenString)

	if err != nil {
		return fmt.Errorf("failed to delete refresh token : %s", err.Error())
	}

	return nil
}

func (d *Database) DeleteAllUserTokens(ctx context.Context, userID int64) error {
	_, err := d.Db.ExecContext(ctx, `
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`, userID)

	if err != nil {
		return fmt.Errorf("failed to delete all user tokens : %s", err.Error())
	}

	return nil
}
