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

	resp, err := http.Post(Config.Slack_Webhook, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to Slack, status code: %d", resp.StatusCode)
	}

	return nil
}

func SendMeetingConclusion(conclusion string, startTime time.Time, duration time.Duration) error {
	message := fmt.Sprintf("Meeting conclusion for %s (Duration: %v):\n%s", startTime.Format("2006-01-02 15:04:05"), duration, conclusion)
	return SendMessage(message)
}
