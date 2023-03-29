package main

import (
	"fmt"
	"os"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

const transcriptFile = "input.txt"

func main() {
	data, err := os.ReadFile(transcriptFile)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", transcriptFile, err)
		return
	}

	transcription := string(data)
	fmt.Printf("Sending:\n%s\n", transcription)

	if err := concluder.SendMessage(transcription); err != nil {
		fmt.Printf("Error sending message to Slack channel: %v\n", err)
		return

	}
	fmt.Println("Message sent to Slack channel successfully.")
}
