// Copyright 2023 Schibsted. Licensed under the terms of the Apache 2.0 license. See LICENSE in the project root.

package concluder

import (
	"bytes"
	"log"
	"math"
	"time"
)

func (a *AudioRecorder) ListenForClapSoundToStopRecording(nClaps int, onRecordingStop func()) {
	audioLength := 1 * time.Second
	loopSleep := 50 * time.Millisecond

	// Clap detection parameters
	energyThreshold := 5000.0              // Adjust this value based on the sensitivity you want
	windowSize := 1024                     // Size of the sliding window used for energy calculation
	clapCooldown := 400 * time.Millisecond // Time to wait before detecting another clap

	clapsDetected := 0
	lastClapDetected := time.Now().Add(-clapCooldown)

	for a.Recording {
		// Create new audio.IntBuffer.
		tailData, err := a.GetRecordedDataTail(audioLength)
		if err != nil {
			log.Println(err)
			time.Sleep(loopSleep)
			continue
		}

		audioBuf, err := newAudioIntBuffer(bytes.NewReader(tailData))
		if err != nil {
			log.Println(err)
			time.Sleep(loopSleep)
			continue
		}

		// Analyze audioBuf to see if it contains a clap sound.
		numSamples := len(audioBuf.Data)
		for i := 0; i < numSamples-windowSize; i++ {
			energy := 0.0
			for j := 0; j < windowSize; j++ {
				energy += math.Abs(float64(audioBuf.Data[i+j]))
			}
			energy /= float64(windowSize)

			if energy > energyThreshold && time.Since(lastClapDetected) > clapCooldown {
				log.Println("GOT A CLAP SOUND")
				clapsDetected++
				lastClapDetected = time.Now()

				if clapsDetected >= nClaps {
					a.StopRecording()
					if onRecordingStop != nil {
						onRecordingStop()
					}
					return
				}
			}
		}

		time.Sleep(loopSleep)
	}
}
