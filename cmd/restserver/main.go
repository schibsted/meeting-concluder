package main

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
	"log"
	"os"
)

func main() {
	ar := concluder.NewAudioRecorder()
	defer ar.Done()

	err := ar.UserSelectsTheInputDevice()
	if err != nil {
		log.Fatal(err)
	}

	audioHandlers := NewAudioHandlers(ar)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", loginHandler)

	// JWT middleware
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	})

	e.GET("/user", getUser, jwtMiddleware)

	e.POST("/start", audioHandlers.startRecordingHandler, jwtMiddleware)
	e.POST("/stop", audioHandlers.stopRecordingHandler, jwtMiddleware)

	addr := os.Getenv("HOST")
	if addr == "" {
		addr = ":3000"
	}
	log.Printf("Starting server on %s...\n", addr)

	if err := e.Start(addr); err != nil {
		log.Fatal(err)
	}
}
