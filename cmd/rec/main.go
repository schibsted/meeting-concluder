package main

import (
	"fmt"
	"os"
	"time"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

func main() {
	// Initialize the AudioRecorder from the concluder package
	audioRecorder := concluder.NewAudioRecorder()

	// Start recording
	const userSelectsDevice = true

	device, err := audioRecorder.UserSelectsTheInputDevice()
	if err != nil {
		fmt.Println("Error letting the user select an input device")
		os.Exit(1)
	}

	audioRecorder.StartRecording(device)

	// Record for 5 seconds
	fmt.Println("Recording for 5 seconds...")
	time.Sleep(5 * time.Second)

	// Stop recording
	audioRecorder.StopRecording()

	// Save the recorded data to a .wav file
	if err := audioRecorder.SaveWav("output.wav"); err != nil {
		fmt.Println("Error saving .wav file:", err)
		os.Exit(1)
	}

	fmt.Println("Audio saved to output.wav")
}
