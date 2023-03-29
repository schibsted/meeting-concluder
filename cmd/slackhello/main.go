package main

import (
	"fmt"
	"log"
	"os"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

func main() {
	slackToken := os.Getenv("SLACK_API_KEY")
	slackChannel := os.Getenv("SLACK_CHANNEL")

	if slackToken == "" || slackChannel == "" {
		log.Fatal("SLACK_API_KEY and SLACK_CHANNEL environment variables must be set.")
	}

	slackClient := concluder.NewSlackClient(slackToken, slackChannel)

	err := slackClient.SendMessage(slackChannel, "hello")
	if err != nil {
		fmt.Printf("Error sending message to Slack channel: %v\n", err)
		return
	}

	fmt.Println("Message sent to Slack channel successfully.")
}
