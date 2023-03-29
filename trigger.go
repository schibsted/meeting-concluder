package concluder

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// THIS FILE IS A WORK IN PROGRESS
//
// Calling the whisper API often enough to be able to pick up trigger words at
// the end of the audio recording is probably not a good idea.

func (a *AudioRecorder) ListenForTriggerWordToStopRecording() error {
	const triggerWavFilename = "/tmp/trigger.wav"
	const triggerWord = "captain"
	audioLength := 3 * time.Second
	coolOff := 10 * time.Second
	loopSleep := audioLength + coolOff
	for a.recording {
		if err := a.SaveTailToWav(3*time.Second, triggerWavFilename); err != nil {
			return fmt.Errorf("error saving %s file: %v", triggerWavFilename, err)
		}
		defer os.Remove("/tmp/trigger.wav")

		transcript, err := TranscribeAudio("/tmp/trigger.wav")
		if err != nil {
			return fmt.Errorf("error transcribing %s: %v", triggerWavFilename, err)
		}
		if strings.Contains(transcript, triggerWord) {
			fmt.Printf("Got trigger word %s\n", triggerWord)
			a.StopRecording()
			return nil
		}
		time.Sleep(loopSleep)
	}
	return nil

}
