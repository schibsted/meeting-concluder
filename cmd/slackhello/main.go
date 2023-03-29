package main

import (
	"fmt"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

func main() {
	slackWebhook := concluder.NewSlackWebhook()

	err := slackWebhook.SendMessage("hello")
	if err != nil {
		fmt.Printf("Error sending message to Slack channel: %v\n", err)
		return
	}

	fmt.Println("Message sent to Slack channel successfully.")
}
