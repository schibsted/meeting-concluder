package concluder

import (
	"fmt"
	"time"

	"github.com/slack-go/slack"
)

type SlackClient struct {
	client  *slack.Client
	channel string
}

func NewSlackClient(config *Config) *SlackClient {
	return &SlackClient{
		client:  slack.New(config.SlackToken),
		channel: config.SlackChannel,
	}
}

func (sc *SlackClient) SendMessage(channel, message string) error {
	if channel == "" {
		channel = sc.channel
	}

	_, _, err := sc.client.PostMessage(channel, slack.MsgOptionText(message, false))
	if err != nil {
		return err
	}

	return nil
}

func (sc *SlackClient) SendMeetingConclusion(conclusion string, startTime time.Time, duration time.Duration) error {
	message := fmt.Sprintf("Meeting conclusion for %s (Duration: %v):\n%s", startTime.Format("2006-01-02 15:04:05"), duration, conclusion)
	return sc.SendMessage(sc.channel, message)
}

func (sc *SlackClient) UpdateConfig(token, channel string) {
	if token != "" {
		sc.client = slack.New(token)
	}
	if channel != "" {
		sc.channel = channel
	}
}
