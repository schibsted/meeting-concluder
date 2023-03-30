package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

var wavFileName = "output.wav"
var mp4FileName = "output.mp4"

func registerRecordingHandlers(e *echo.Echo, audioRecorder *concluder.AudioRecorder) {
	e.GET("/devices", getDevices(audioRecorder))
	e.POST("/select-device/:index", selectDevice(audioRecorder))
	e.POST("/record", startRecording(audioRecorder))
	e.POST("/stop", stopRecording(audioRecorder))
	e.GET("/conclusion", getConclusion())
}

func getDevices(audioRecorder *concluder.AudioRecorder) echo.HandlerFunc {
	return func(c echo.Context) error {
		devices, err := audioRecorder.InputDevices()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get input devices"})
		}

		deviceList := make([]map[string]interface{}, len(devices))
		for i, device := range devices {
			deviceList[i] = map[string]interface{}{
				"index": i,
				"name":  device.Name,
			}
		}

		return c.JSON(http.StatusOK, deviceList)
	}
}

func selectDevice(audioRecorder *concluder.AudioRecorder) echo.HandlerFunc {
	return func(c echo.Context) error {
		index := c.Param("index")
		devices, err := audioRecorder.InputDevices()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error getting input devices"})
		}
		selectedDevice, err := strconv.Atoi(index)
		if err != nil || selectedDevice < 0 || selectedDevice >= len(devices) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid device index"})
		}
		audioRecorder.SetSelectedDevice(devices[selectedDevice])
		return c.JSON(http.StatusOK, map[string]string{"message": "Device selected successfully"})
	}
}

func startRecording(audioRecorder *concluder.AudioRecorder) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if recording is already in progress
		if audioRecorder.IsRecording() {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Recording is already in progress"})
		}

		// Parse the recording duration from the request
		durationStr := c.QueryParam("duration")
		if durationStr == "" {
			durationStr = "3600s"
		}
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid duration"})
		}

		// Start the recording in a separate goroutine
		go func() {
			var err error
			if err := audioRecorder.RecordAudio(wavFileName, duration, nClapDetection); err != nil {
				c.Logger().Errorf("Error recording audio to %s: %v", wavFileName, err)
				return
			}

			c.Logger().Info("Transcribing and concluding")
			conclusion, err = audioRecorder.TranscribeConvertConclude(wavFileName, mp4FileName, true, true)
			if err != nil {
				c.Logger().Errorf("Error generating conclusion: %v", err)
				return
			}

			c.Logger().Infof("Generated conclusion: %s", conclusion)
		}()

		return c.JSON(http.StatusOK, map[string]string{"message": "Recording started"})
	}
}

func getConclusion() echo.HandlerFunc {
	return func(c echo.Context) error {
		if conclusion == "" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Conclusion not available"})
		}

		return c.JSON(http.StatusOK, map[string]string{"conclusion": conclusion})
	}
}

func stopRecording(audioRecorder *concluder.AudioRecorder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !audioRecorder.IsRecording() {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Recording is not in progress"})
		}

		audioRecorder.StopRecording()

		conclusion, err := audioRecorder.TranscribeConvertConclude(wavFileName, mp4FileName, true, true)
		if err != nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Could not conclude"})
		}

		c.Logger().Infof("Generated conclusion: %s", conclusion)

		return c.JSON(http.StatusOK, map[string]string{"message": "Recording stopped by user", "conclusion": conclusion})
	}
}
