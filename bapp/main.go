package main

import (
	"bapp/db"
	"bapp/helpers"
	"bapp/routes"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	dbInstance, err := db.DB()
	if err != nil {
		panic(err)
	}
	handler := &routes.HandlerClient{
		DBClient: db.NewClient(dbInstance),
	}

	sqlDB, err := dbInstance.DB()
	if err != nil {
		fmt.Println("Error getting underlying DB:", err)
		return
	}
	defer sqlDB.Close()

	// Define routes and handlers here
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/users", handler.RegisterUser)
	e.POST("/login", handler.LoginUser)
	e.GET("/userconfig", helpers.ValidateJWT(handler.GetUserConfig))
	e.PATCH("/userconfig", helpers.ValidateJWT(handler.UpdateUserConfig))

	// Start the Echo server
	e.Start(":8080")

}
