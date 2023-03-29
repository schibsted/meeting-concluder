package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

func main() {
	// Initialize the AudioRecorder from the concluder package
	audioRecorder := concluder.NewAudioRecorder()

	// Let the user select an input device through a CLI menu
	if err := audioRecorder.UserSelectsTheInputDevice(); err != nil {
		fmt.Println("Error letting the user select an input device")
		os.Exit(1)
	}

	// Start recording
	audioRecorder.StartRecording()

	// Set up a channel to listen for interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Record for 10 seconds or until Ctrl-C is pressed
	fmt.Println("Recording for 10 seconds or until Ctrl-C is pressed...")
	select {
	case <-time.After(10 * time.Second):
	case <-signalChan:
		fmt.Println("\nCtrl-C received, stopping recording.")
	}

	// Stop recording
	audioRecorder.StopRecording()

	// Save the recorded data to a .wav file
	if err := audioRecorder.SaveWav("output.wav"); err != nil {
		fmt.Println("Error saving .wav file:", err)
		os.Exit(1)
	}

	fmt.Println("Audio saved to output.wav")
}
