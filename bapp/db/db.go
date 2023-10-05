package db

import (
	"encoding/json"
	"fmt"
	"os"

	model "bapp/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

func DB() (*gorm.DB, error) {
	var config Config
	plan, err := os.ReadFile("db/config.json")
	if err != nil {
		log.Error().Err(err).Msg("Error reading config.json")
		return nil, err
	}
	readErr := json.Unmarshal(plan, &config)
	if readErr != nil {
		log.Error().Err(readErr).Msg("Error unmarshalling config.json")
		return nil, readErr
	}

	info := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := gorm.Open(postgres.Open(info), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to database")
		return nil, err
	}

	// we need to do auto migration for our models
	migrateErr := db.AutoMigrate(
		&model.User{},
		&model.UserConfig{},
	)
	if migrateErr != nil {
		log.Error().Err(migrateErr).Msg("Error migrating models")
		return nil, migrateErr
	}

	return db, nil
}
