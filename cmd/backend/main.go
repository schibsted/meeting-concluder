package main

import (
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
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	audioHandlers := NewAudioHandlers()
	userHandlers := NewUserHandlers()
	userHandlers.Configure(e)

	e.GET("/devices", audioHandlers.listDevicesHandler)
	e.POST("/selectdevice", audioHandlers.selectDeviceHandler)
	e.POST("/start", audioHandlers.startRecordingHandler)
	e.POST("/stop", audioHandlers.stopRecordingHandler)
	e.POST("/process", audioHandlers.transcribeConvertConcludeHandler, userHandlers.JWTMiddleware())

	addr := ":3000"
	e.Logger.Fatal(e.Start(addr))
}
