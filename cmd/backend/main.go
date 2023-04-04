// Copyright 2023 Schibsted. Licensed under the terms of the Apache 2.0 license. See LICENSE in the project root.

package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

const nClapDetection = 2 // number of claps detected for the recording to stop, use 0 to disable

var (
	conclusion  string
	maxDuration time.Duration
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize the AudioRecorder from the concluder package
	audioRecorder := concluder.NewAudioRecorder()
	defer audioRecorder.Terminate()

	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	e.GET("/:filename", func(c echo.Context) error {
		filename := c.Param("filename")
		return c.File("./static/" + filename)
	})

	registerRecordingHandlers(e, audioRecorder)

	e.Logger.Fatal(e.Start(":4000"))
}
