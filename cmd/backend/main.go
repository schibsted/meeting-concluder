package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

const nClapDetection = 2 // number of claps detected for the recording to stop, use 0 to disable

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

	e.GET("/static/:filename", func(c echo.Context) error {
		filename := c.Param("filename")
		return c.File("./static/" + filename)
	})

	e.GET("/devices", func(c echo.Context) error {
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
	})

	e.POST("/select-device/:index", func(c echo.Context) error {
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
	})

	e.POST("/record", func(c echo.Context) error {
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
			wavFileName, err := audioRecorder.RecordAudio(duration, nClapDetection)
			if err != nil {
				c.Logger().Errorf("Error recording audio: %v", err)
				return
			}
			defer os.Remove(wavFileName)

			c.Logger().Info("Transcribing and concluding")
			conclusion, err := audioRecorder.TranscribeConvertConclude(wavFileName)
			if err != nil {
				c.Logger().Errorf("Error generating conclusion: %v", err)
				return
			}

			c.Logger().Infof("Generated conclusion: %s", conclusion)
		}()

		return c.JSON(http.StatusOK, map[string]string{"message": "Recording started"})
	})

	e.POST("/stop", func(c echo.Context) error {
		if !audioRecorder.IsRecording() {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Recording is not in progress"})
		}

		audioRecorder.StopRecording()
		return c.JSON(http.StatusOK, map[string]string{"message": "Recording stopped"})
	})

	e.Logger.Fatal(e.Start(":3000"))
}
