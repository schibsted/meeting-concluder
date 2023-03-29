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
	audioRecorder.StartRecording(userSelectsDevice)

	// Record for 5 seconds
	fmt.Println("Recording for 5 seconds...")
	time.Sleep(5 * time.Second)

	// Stop recording
	audioRecorder.StopRecording()

	// Save the recorded data to a .wav file
	err := audioRecorder.SaveWav("output.wav")
	if err != nil {
		fmt.Println("Error saving .wav file:", err)
		os.Exit(1)
	}

	fmt.Println("Audio saved to output.wav")
}
