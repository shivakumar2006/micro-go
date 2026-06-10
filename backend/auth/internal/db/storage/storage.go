package db

import (
	"auth/internal/model"
	"context"
)

type Storage interface {
	CreateUser(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	FindUserById(id int64) (*model.User, error)
	SaveRefreshToken(token *model.RefreshToken) error
	FindRefreshToken(token *model.RefreshToken) error
	DeleteRefreshToken(tokenString string) error
	DeleteAllUserTokens(ctx context.Context, userID int64) error
}
