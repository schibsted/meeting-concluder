package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

type AudioHandlers struct {
	ar *concluder.AudioRecorder
}

func NewAudioHandlers() *AudioHandlers {
	ar := concluder.NewAudioRecorder()
	return &AudioHandlers{ar: ar}
}

func (ah *AudioHandlers) listDevicesHandler(c echo.Context) error {
	devices, err := ah.ar.ListDevices()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, devices)
}

func (ah *AudioHandlers) selectDeviceHandler(c echo.Context) error {
	deviceID, err := strconv.Atoi(c.FormValue("device_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid device ID")
	}

	err = ah.ar.SelectDevice(deviceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Device selected"})
}

func (ah *AudioHandlers) startRecordingHandler(c echo.Context) error {
	if ah.ar.IsRecording() {
		return echo.NewHTTPError(http.StatusBadRequest, "Already recording")
	}

	clapDetection := c.FormValue("clap_detection") == "true"
	duration, err := strconv.Atoi(c.FormValue("duration"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid duration")
	}

	go func() {
		filename, err := ah.ar.RecordAudio(clapDetection, time.Duration(duration)*time.Second)
		if err != nil {
			c.Logger().Errorf("Error recording audio: %v", err)
		}
		ah.ar.SetOutputFilename(filename)
	}()

	return c.JSON(http.StatusOK, map[string]string{"message": "Recording started"})
}

func (ah *AudioHandlers) stopRecordingHandler(c echo.Context) error {
	if !ah.ar.IsRecording() {
		return echo.NewHTTPError(http.StatusBadRequest, "Not recording")
	}

	ah.ar.StopRecording()
	return c.JSON(http.StatusOK, map[string]string{"message": "Recording stopped"})
}

func (ah *AudioHandlers) transcribeConvertConcludeHandler(c echo.Context) error {
	outputFilename := ah.ar.GetOutputFilename()
	if outputFilename == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No recording available")
	}

	conclusion, err := ah.ar.TranscribeConvertConcludePost(outputFilename)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"conclusion": conclusion})
}
