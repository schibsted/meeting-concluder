// Copyright 2023 Schibsted. Licensed under the terms of the Apache 2.0 license. See LICENSE in the project root.

package concluder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func SendMessage(message string) error {
	payload := map[string]interface{}{
		"text": message,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(Config.SlackWebhook, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to Slack, status code: %d", resp.StatusCode)
	}

	return nil
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func SendMeetingConclusion(conclusion string, startTime, endTime time.Time) error {
	duration := endTime.Sub(startTime)
	formattedDuration := formatDuration(duration)
	message := fmt.Sprintf("Meeting conclusion for %s (Duration: %s):\n%s", startTime.Format("2006-01-02 15:04:05"), formattedDuration, conclusion)
	return SendMessage(message)
}
