package db

import (
	"bapp/models"

	"github.com/google/uuid"
)

type DBInterface interface {
	CreateUser(user *models.User) error
	UserExists(email string) bool
	GetUserByEmail(email string) (*models.User, error)
	CreateUserConfig(userConfig *models.UserConfig) error
	GetUserConfig(userID uuid.UUID) (*models.UserConfig, error)
	UpdateUserConfig(userConfig *models.UserConfig) error
}
