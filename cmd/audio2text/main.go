package main

import (
	"fmt"
	"os"
	"strings"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

const mp4FilePath = "input.mp4"

func main() {
	transcript, err := concluder.Transcribe(mp4FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error transribing %s: %v\n", mp4FilePath, err)
		return
	}

	transcript += "\n"
	fmt.Printf("%s\n", transcript)

	if err := os.WriteFile("output.txt", []byte(transcript), 0o644); err != nil {
		fmt.Printf("Error writing to output.txt: %v\n", err)
		return
	}
}
