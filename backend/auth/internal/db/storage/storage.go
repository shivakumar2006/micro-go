package storage

import "auth/internal/models"

type Storage interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserById(id int) (*models.User, error)
	SaveRefreshToken(session *models.UserSessions) error
	FindRefreshToken(refreshTokenHash string) (*models.UserSessions, error)
	DeleteRefreshToken(refreshTokenHash string) error
	DeleteAllUserToken(userId int) error
}
