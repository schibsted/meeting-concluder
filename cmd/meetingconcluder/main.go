package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

const nClapDetection = 2 // number of claps detected for the recording to stop, use 0 to disable

func main() {
	// Initialize the AudioRecorder from the concluder package
	audioRecorder := concluder.NewAudioRecorder()
	defer audioRecorder.Terminate()

	// Let the user select an input device through a CLI menu
	if err := audioRecorder.UserSelectsTheInputDevice(); err != nil {
		fmt.Println("Error letting the user select an input device")
		os.Exit(1)
	}

	startTime := time.Now()

	// Create a temporary file to store the recorded audio
	tmpfile, err := ioutil.TempFile("", "audio-*.wav")
	if err != nil {
		fmt.Printf("Error creating temporary file: %v", err)
		os.Exit(1)
	}
	defer os.Remove(tmpfile.Name())

	// Record audio to the temporary file
	fmt.Printf("Recording audio. To stop before the specified max duration, press ctrl-c or clap %d time(s)...\n", nClapDetection)

	// Create a channel to listen for ctrl-c interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Wait for the interrupt signal
		<-sigCh

		// Stop recording
		audioRecorder.StopRecording()
	}()

	if err := audioRecorder.RecordToFile(tmpfile.Name(), 1*time.Hour, nClapDetection); err != nil {
		fmt.Printf("Error recording audio to file: %v", err)
		os.Exit(1)
	}

	conclusion, err := audioRecorder.TranscribeConvertConclude(tmpfile.Name())
	if err != nil {
		fmt.Printf("Could not conclude: %v", err)
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	concluder.SendMeetingConclusion(conclusion, startTime, duration)
}
