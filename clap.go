package concluder

import (
	"bytes"
	"log"
	"math"
	"time"
)

func (a *AudioRecorder) ListenForTripleClapSoundToStopRecording() {
	audioLength := 1 * time.Second
	loopSleep := 50 * time.Millisecond

	// Clap detection parameters
	energyThreshold := 4000.0 // Adjust this value based on the sensitivity you want
	windowSize := 1024        // Size of the sliding window used for energy calculation
	clapCount := 0            // Number of claps detected

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

			if energy > energyThreshold {
				clapCount++
				if clapCount == 3 {
					log.Println("GOT TRIPLE CLAPS SOUND")
					a.stopRecording()
					return
				}
			} else {
				clapCount = 0
			}
		}

		time.Sleep(loopSleep)
	}
}

func (a *AudioRecorder) ListenForClapSoundToStopRecording() {
	audioLength := 1 * time.Second
	loopSleep := 50 * time.Millisecond

	// Clap detection parameters
	energyThreshold := 4000.0 // Adjust this value based on the sensitivity you want
	windowSize := 1024        // Size of the sliding window used for energy calculation

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

			if energy > energyThreshold {
				log.Println("GOT A CLAP SOUND")
				a.stopRecording()
				return
			}
		}

		time.Sleep(loopSleep)
	}
}
