package repository

import (
	"auth/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepository struct {
	Db *sql.DB
}

func NewAuthRepository(db *sql.DB) (*AuthRepository, error) {
	return &AuthRepository{
		Db: db,
	}, nil
}

func (a *AuthRepository) CreateUser(user *models.User) error {
	ctx, cancel := NewContext()
	defer cancel()

	err := a.Db.QueryRowContext(ctx, `
		INSERT INTO users(name, email, password, role_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email
	`, user.Name, user.Email, user.Password, user.RoleID).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("user already exists : %w", err)
			}
		}
		return fmt.Errorf("failed to create user : %w", err)
	}

	return nil
}

func (a *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	ctx, cancel := NewContext()
	defer cancel()

	var user models.User

	err := a.Db.QueryRowContext(ctx, `
		SELECT u.id, u.name, u.email, u.password, u.role_id, r.role_name
		FROM users u
		JOIN roles r 
		ON u.role_id = r.id
		WHERE u.email = $1
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.RoleName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find user by email : %w", err)
	}

	return &user, nil
}

func (a *AuthRepository) FindUserById(id int) (*models.User, error) {
	ctx, cancel := NewContext()
	defer cancel()

	var user models.User

	err := a.Db.QueryRowContext(ctx, `
		SELECT u.id, u.name, u.email, u.password, u.role_id, r.role_name
		FROM users u
		JOIN roles r 
		ON u.role_id = r.id
		WHERE u.id = $1
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.RoleName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found : %w", err)
		}
		return nil, fmt.Errorf("failed to find user by id : %w", err)
	}

	return &user, nil
}

func (a *AuthRepository) SaveRefreshToken(session *models.UserSessions) error {
	ctx, cancel := NewContext()
	defer cancel()

	_, err := a.Db.ExecContext(ctx, `
		INSERT INTO user_sessions(user_id, refresh_token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`, session.UserID, session.RefreshTokenHash, session.ExpiresAt, session.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("refresh token already exists : %w", err)
			}
		}
		return fmt.Errorf("failed to save refresh token : %w", err)
	}

	return nil
}

func (a *AuthRepository) FindRefreshToken(refreshTokenHash string) (*models.UserSessions, error) {
	ctx, cancel := NewContext()
	defer cancel()

	var session models.UserSessions

	err := a.Db.QueryRowContext(ctx, `
		SELECT id, user_id, refresh_token_hash, expires_at, created_at
		FROM user_sessions
		WHERE refresh_token_hash = $1
	`, refreshTokenHash).Scan(&session.ID, &session.UserID, &session.RefreshTokenHash, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("refresh token not found : %w", err)
		}
		return nil, fmt.Errorf("failed to find refresh token : %w", err)
	}

	return &session, nil
}

func (a *AuthRepository) DeleteRefreshToken(refreshTokenHash string) error {
	ctx, cancel := NewContext()
	defer cancel()

	_, err := a.Db.ExecContext(ctx, `
		DELETE FROM user_sessions
		WHERE refresh_token_hash = $1
	`, refreshTokenHash)

	if err != nil {
		return fmt.Errorf("failed to delete refresh token : %w", err)
	}

	return nil
}

func (a *AuthRepository) DeleteAllUserToken(userId int) error {
	ctx, cancel := NewContext()
	defer cancel()

	_, err := a.Db.ExecContext(ctx, `
		DELETE FROM user_sessions
		WHERE user_id = $1
	`, userId)

	if err != nil {
		return fmt.Errorf("failed to delete all users token")
	}

	return nil
}

func (a *AuthRepository) GetRoleByName(roleName string) (int, error) {
	ctx, cancel := NewContext()
	defer cancel()

	var role models.Role

	err := a.Db.QueryRowContext(ctx, `
		SELECT id
		FROM roles
		WHERE role_name = $1
	`, roleName).Scan(&role.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("role not found : %w", err)
		}
		return 0, fmt.Errorf("failed to get role by name : %w", err)
	}

	return role.ID, nil
}

func (a *AuthRepository) AdminExist() (bool, error) {
	ctx, cancel := NewContext()
	defer cancel()

	var exists bool

	err := a.Db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 
			FROM users u 
			JOIN roles r ON u.role_id = r.id
			WHERE r.role_name = $1
		);
	`, "admin").Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check admin exists : %w", err)
	}

	return exists, nil
}

func NewContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
