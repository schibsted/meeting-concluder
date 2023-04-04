// Copyright 2023 Schibsted. Licensed under the terms of the Apache 2.0 license. See LICENSE in the project root.

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
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

	convertMutex sync.Mutex
)

func registerRecordingHandlers(e *echo.Echo, audioRecorder *concluder.AudioRecorder) {
	e.GET("/devices", getDevices(audioRecorder))
	e.POST("/select-device/:index", selectDevice(audioRecorder))
	e.POST("/record", startRecording(audioRecorder))
	e.POST("/stop", stopRecording(audioRecorder))
	e.GET("/conclusion", getConclusion())
	e.POST("/post-to-slack", postToSlack(audioRecorder))
	e.POST("/stopped-by-clapping", stoppedByClapping())
	e.GET("/clap-stop-event", clapStopEvent())
	e.DELETE("/conclusion", clearConclusion())
}

var clapStopChan = make(chan struct{})

func clapStopEvent() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")

		select {
		case <-clapStopChan:
			c.Response().Write([]byte("data: Recording stopped by clapping\n\n"))
			c.Response().Flush()
		case <-c.Request().Context().Done():
		}
		return nil
	}
}

func stoppedByClapping() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Recording stopped by clapping"})
	}
}

func sendClapStopSignal() {
	clapStopChan <- struct{}{}
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

		convertMutex.Lock()
		conclusion = ""
		convertMutex.Unlock()

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
			convertMutex.Lock()
			if err := audioRecorder.RecordAudio(wavFileName, duration, nClapDetection, sendClapStopSignal); err != nil {
				c.Logger().Errorf("Error recording audio to %s: %v", wavFileName, err)
				convertMutex.Unlock()
				return
			}
			convertMutex.Unlock()

			stopTime = time.Now()

			c.Logger().Info("Transcribing and concluding")
			convertMutex.Lock()
			if conclusion == "" {
				conclusion, err = concluder.TranscribeConvertConclude(wavFileName, mp4FileName, true, true)
				if err != nil {
					c.Logger().Errorf("Error generating conclusion: %v", err)
					convertMutex.Unlock()
					return
				}
			}
			convertMutex.Unlock()

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

func clearConclusion() echo.HandlerFunc {
	return func(c echo.Context) error {
		convertMutex.Lock()
		conclusion = ""
		convertMutex.Unlock()
		return c.JSON(http.StatusOK, map[string]string{"message": "Conclusion cleared"})
	}
}

func stopRecording(audioRecorder *concluder.AudioRecorder) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !audioRecorder.IsRecording() {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Recording is not in progress"})
		}

		audioRecorder.StopRecording()
		stopTime = time.Now()

		var err error
		convertMutex.Lock()
		if conclusion == "" {
			conclusion, err = concluder.TranscribeConvertConclude(wavFileName, mp4FileName, true, true)
			if err != nil {
				convertMutex.Unlock()
				return c.JSON(http.StatusConflict, map[string]string{"error": "Could not conclude"})
			}
		}
		convertMutex.Unlock()

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
