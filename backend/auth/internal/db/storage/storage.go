package db

import "auth/internal/model"

type Storage interface {
	CreateUser(user *model.User) error
}
