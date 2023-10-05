package helpers

import (
	"bapp/models"
	"log"
	"net/http"
	"reflect"
	"strings"

	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var secret = "victoria"

func GenerateJWT(user models.User) (string, error) {
	log.Println(user.Email)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   user.Email,
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token != "" {
			// we need to remove the Bearer prefix from the token
			parts := strings.Split(token, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
			}
			token = parts[1]
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					log.Println("Error parsing token1")
					return "", nil
				}
				return []byte(secret), nil
			})
			if err != nil {
				log.Printf("Error parsing token: %T - %s\n", err, err) // Print error details
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			if parsedToken.Valid {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
		return echo.ErrUnauthorized
	}
}

func ComparePasswords(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func UpdateStructValues(struct1, struct2 interface{}) {
	value1 := reflect.ValueOf(struct1).Elem()
	value2 := reflect.ValueOf(struct2).Elem()

	for i := 0; i < value1.NumField(); i++ {
		field1 := value1.Field(i)
		field2 := value2.Field(i)

		// Compare values and update if different
		if !reflect.DeepEqual(field1.Interface(), field2.Interface()) {
			field1.Set(field2)
		}
	}
}
