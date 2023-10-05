package routes

import (
	"bapp/db"
	"bapp/helpers"
	"bapp/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type HandlerClient struct {
	DBClient db.DBInterface
}

type Handler interface {
	RegisterUser(c echo.Context) error
	LoginUser(c echo.Context) error
}

func (h *HandlerClient) RegisterUser(c echo.Context) error {
	var user models.User
	bindErr := c.Bind(&user)
	if bindErr != nil {
		return bindErr
	}
	// check if the user already exists
	if h.DBClient.UserExists(user.Email) {
		return c.JSON(400, "User already exists")
	}
	createErr := h.DBClient.CreateUser(&user)
	if createErr != nil {
		return createErr
	}
	// create a user config for the user
	userConfig := models.UserConfig{
		UserID:  user.ID,
		Balance: 1000,
	}
	userConfigErr := h.DBClient.CreateUserConfig(&userConfig)
	if userConfigErr != nil {
		return userConfigErr
	}
	return c.JSON(200, user)

}

// a function to login a user
func (h *HandlerClient) LoginUser(c echo.Context) error {
	var creds models.UserLogin
	bindErr := c.Bind(&creds)
	if bindErr != nil {
		return bindErr
	}
	// check if the user already exists
	if !h.DBClient.UserExists(creds.Email) {
		return c.JSON(400, "User does not exist")
	}
	// get the user from the database
	user, userErr := h.DBClient.GetUserByEmail(creds.Email)
	if userErr != nil {
		return userErr
	}
	log.Info().Msg("User found")
	// check if the password is correct
	if !helpers.ComparePasswords(user.PasswordHash, creds.Password) {
		return c.JSON(401, "Incorrect password")
	}

	jwtToken, jwtErr := helpers.GenerateJWT(*user)
	if jwtErr != nil {
		return jwtErr
	}
	return c.JSON(200, jwtToken)
}

// a function to get a user config
func (h *HandlerClient) GetUserConfig(c echo.Context) error {
	userID := c.QueryParam("userId")
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	userConfig, userConfigErr := h.DBClient.GetUserConfig(parsedID)
	if userConfigErr != nil {
		return userConfigErr
	}
	return c.JSON(200, userConfig)
}

// a function to update a user config
func (h *HandlerClient) UpdateUserConfig(c echo.Context) error {
	var userConfig models.UserConfig
	bindErr := c.Bind(&userConfig)
	if bindErr != nil {
		fmt.Println("Error binding user config:", bindErr)
		return bindErr
	}
	// update the user config
	newConfigErr := h.DBClient.UpdateUserConfig(&userConfig)
	if newConfigErr != nil {
		return newConfigErr
	}
	return c.JSON(200, "User config updated")
}
