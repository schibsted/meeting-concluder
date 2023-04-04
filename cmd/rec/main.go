// Copyright 2023 Schibsted. Licensed under the terms of the Apache 2.0 license. See LICENSE in the project root.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

const (
	wavFilename    = "output.wav"
	nClapDetection = 2 // detect N number of claps before stopping. 0 to disable
)

func main() {
	// Initialize the AudioRecorder from the concluder package
	audioRecorder := concluder.NewAudioRecorder()
	defer audioRecorder.Terminate()

	// Let the user select an input device through a CLI menu
	if err := audioRecorder.UserSelectsTheInputDevice(); err != nil {
		fmt.Println("Error letting the user select an input device")
		os.Exit(1)
	}

	// Record audio to the specified file
	fmt.Printf("Recording audio. To stop before the specified max duration, press ctrl-c or clap %d time(s)...\n", nClapDetection)

	// Create a channel to listen for ctrl-c interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Wait for the interrupt signal
		<-sigCh

		// Close the StopRecordingCh channel to stop recording
		close(audioRecorder.StopRecordingCh)
	}()

	if err := audioRecorder.RecordToFile(wavFilename, 1*time.Hour, nClapDetection, nil); err != nil {
		fmt.Printf("Error recording audio to file: %v", err)
		os.Exit(1)
	}
}
