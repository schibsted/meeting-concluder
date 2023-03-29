package main

import (
	"fmt"
	"os"
	"time"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

func main() {
	const wavFilename = "output.wav"

	// Initialize the AudioRecorder from the concluder package
	audioRecorder := concluder.NewAudioRecorder()
	defer audioRecorder.Done()

	// Let the user select an input device through a CLI menu
	if err := audioRecorder.UserSelectsTheInputDevice(); err != nil {
		fmt.Println("Error letting the user select an input device")
		os.Exit(1)
	}

	// Record audio to the specified file
	if err := audioRecorder.RecordToFile(wavFilename, 10*time.Second); err != nil {
		fmt.Printf("Error recording audio to file: %v", err)
		os.Exit(1)
	}
}
