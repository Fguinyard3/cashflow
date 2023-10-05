package db

import (
	model "bapp/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DBClient struct {
	gorm *gorm.DB
}

func NewClient(gormDB *gorm.DB) *DBClient {
	return &DBClient{gorm: gormDB}
}

func (c *DBClient) CreateUser(user *model.User) error {
	log.Info().Msg("Creating user")
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return err
	}
	user.PasswordHash = string(bcryptPassword)
	return c.gorm.Create(user).Error
}

// A function to create a userconfig
func (c *DBClient) CreateUserConfig(userConfig *model.UserConfig) error {
	log.Info().Msg("Creating user config")
	return c.gorm.Create(userConfig).Error

}

// A function to get a userconfig
func (c *DBClient) GetUserConfig(userID uuid.UUID) (*model.UserConfig, error) {
	log.Info().Msg("Getting user config")
	var userConfig model.UserConfig
	err := c.gorm.Where("user_id = ?", userID).First(&userConfig).Error
	if err != nil {
		return nil, err
	}
	return &userConfig, nil
}

// A function to update a userconfig
func (c *DBClient) UpdateUserConfig(userConfig *model.UserConfig) error {
	log.Info().Msg("Updating user config")
	// first we need to get the user config from the database
	var userConfigFromDB model.UserConfig
	err := c.gorm.Where("user_id = ?", userConfig.UserID.String()).First(&userConfigFromDB).Error
	if err != nil {
		log.Error().Err(err).Msg("Error getting user config from database")
		return err
	}
	fmt.Println("ping")
	// update the user config
	userConfigFromDB.Balance = userConfig.Balance
	userConfigFromDB.Income = userConfig.Income
	userConfigFromDB.Paydates = userConfig.Paydates
	userConfigFromDB.Bills = userConfig.Bills
	fmt.Println("ping")
	fmt.Println(userConfigFromDB)
	saveErr := c.gorm.Save(&userConfigFromDB).Error
	if saveErr != nil {
		return saveErr
	}
	return nil
}

// a function to check if a user exists in the database
func (c *DBClient) UserExists(email string) bool {
	log.Info().Msg("Checking if user exists")
	var user model.User
	result := c.gorm.Where("email = ?", email).First(&user)
	return result.RowsAffected != 0
}

// a function to get a user by email
func (c *DBClient) GetUserByEmail(email string) (*model.User, error) {
	log.Info().Msg("Getting user by email")
	var user model.User
	err := c.gorm.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
