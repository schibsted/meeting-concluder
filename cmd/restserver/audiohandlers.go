package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

type AudioHandlers struct {
	ar *concluder.AudioRecorder
}

func NewAudioHandlers(ar *concluder.AudioRecorder) *AudioHandlers {
	return &AudioHandlers{ar: ar}
}

func (ah *AudioHandlers) startRecordingHandler(c echo.Context) error {
	clapDetection := c.QueryParam("clapDetection")
	duration := c.QueryParam("duration")

	clapDetectEnabled := false
	if clapDetection == "true" {
		clapDetectEnabled = true
	}

	durationSeconds, err := strconv.Atoi(duration)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid duration provided",
		})
	}

	ah.ar.RecordToFile(ah.generateFilename(), time.Duration(durationSeconds)*time.Second, clapDetectEnabled)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Recording started",
	})
}

func (ah *AudioHandlers) stopRecordingHandler(c echo.Context) error {
	ah.ar.StopRecording()
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Recording stopped",
	})
}

func (ah *AudioHandlers) generateFilename() string {
	return fmt.Sprintf("output-%s.wav", time.Now().Format("2006-01-02-15-04-05"))
}
