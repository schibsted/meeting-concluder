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

func NewSlackClient(slackToken, slackChannel string) *SlackClient {
	return &SlackClient{
		client:  slack.New(slackToken),
		channel: slackChannel,
	}
}

func (sc *SlackClient) SendMessage(channel, message string) error {
	if channel == "" {
		channel = sc.channel
	}
	_, _, err := sc.client.PostMessage(channel, slack.MsgOptionText(message, false))
	return err
}

func (sc *SlackClient) SendMeetingConclusion(conclusion string, startTime time.Time, duration time.Duration) error {
	message := fmt.Sprintf("Meeting conclusion for %s (Duration: %v):\n%s", startTime.Format("2006-01-02 15:04:05"), duration, conclusion)
	return sc.SendMessage(sc.channel, message)
}
