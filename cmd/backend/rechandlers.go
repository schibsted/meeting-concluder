package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

func generateFileNameWithTimestamp(prefix, extension string) string {
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	return fmt.Sprintf("%s-%s.%s", prefix, timestamp, extension)
}

var (
	wavFileName = "output.wav" // generateFileNameWithTimestamp("recording", "wav")
	mp4FileName = "output.mp4" // generateFileNameWithTimestamp("recording", "mp4")

	startTime = time.Now()
	stopTime  = time.Now()
)

func registerRecordingHandlers(e *echo.Echo, audioRecorder *concluder.AudioRecorder) {
	e.GET("/devices", getDevices(audioRecorder))
	e.POST("/select-device/:index", selectDevice(audioRecorder))
	e.POST("/record", startRecording(audioRecorder))
	e.POST("/stop", stopRecording(audioRecorder))
	e.GET("/conclusion", getConclusion())
	e.POST("/post-to-slack", postToSlack(audioRecorder))
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

		startTime = time.Now()

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
			if err := audioRecorder.RecordAudio(wavFileName, duration, nClapDetection, func() {
				// Send a POST request to this server to stop the recording
				url := fmt.Sprintf("http://%s/stop", c.Request().Host)
				http.Post(url, "application/json", nil)
			}); err != nil {
				c.Logger().Errorf("Error recording audio to %s: %v", wavFileName, err)
				return
			}
			stopTime = time.Now()

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
		stopTime = time.Now()

		conclusion, err := audioRecorder.TranscribeConvertConclude(wavFileName, mp4FileName, true, true)
		if err != nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Could not conclude"})
		}

		c.Logger().Infof("Generated conclusion: %s", conclusion)

		return c.JSON(http.StatusOK, map[string]string{"message": "Recording stopped by user", "conclusion": conclusion})
	}
}

func postToSlack(audioRecorder *concluder.AudioRecorder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if startTime.IsZero() || stopTime.IsZero() {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Recording start or end time is not available"})
		}
		if err := concluder.SendMeetingConclusion(conclusion, startTime, stopTime); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error sending conclusion to Slack"})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Conclusion posted to Slack"})
	}
}
