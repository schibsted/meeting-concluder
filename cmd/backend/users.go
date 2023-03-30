package main

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

func loginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Replace the following lines with your authentication logic
	if username == "testuser" && password == "testpassword" {
		token, err := createJWTToken(username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Error generating JWT token",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": token,
		})
	}

	return echo.ErrUnauthorized
}

func createJWTToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func getUser(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return errors.New("JWT token missing or invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("failed to cast claims as jwt.MapClaims")
	}
	return c.JSON(http.StatusOK, claims)
}
