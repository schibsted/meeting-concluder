package concluder

import (
	"math"
	"math/cmplx"
	"time"

	"github.com/mjibson/go-dsp/fft"
)

// Parameters for clap detection
const (
	windowSize    = 1024
	clapThreshold = 500.0
	maxInterval   = 0.8
)

// Listen for a specific number of claps within a certain time frame
func ListenForNClapsToStopRecording(a *AudioRecorder, n int) {
	loopSleep := 50 * time.Millisecond

	// Buffer for audio data
	audioBuf := make([]float32, windowSize)
	outputBuf := make([]float32, windowSize)

	// Magnitude buffer for FFT
	magnitude := make([]float64, windowSize)

	// Number of claps detected so far
	clapsDetected := 0

	// Timestamp of last clap detected
	lastClap := time.Now()

	for a.Recording {

		// Read audio data into buffer
		a.captureAudio(audioBuf, outputBuf)

		// Take FFT of audio buffer
		for i, sample := range audioBuf {
			magnitude[i] = float64(sample)
		}
		fft.FFTReal(magnitude)

		// Find index of peak magnitude
		peakIndex := 0
		for i := 1; i < windowSize/2; i++ {
			if math.Abs(magnitude[i]) > math.Abs(magnitude[peakIndex]) {
				peakIndex = i
			}
		}

		// Check if peak magnitude is above threshold
		if cmplx.Abs(complex(magnitude[peakIndex], 0)) > clapThreshold {
			// Check if it's been less than maxInterval since the last clap
			now := time.Now()
			interval := now.Sub(lastClap)
			if interval.Seconds() < maxInterval {
				clapsDetected++
			} else {
				clapsDetected = 1
			}

			lastClap = now

			// Check if we've reached the desired number of claps
			if clapsDetected >= n {
				a.StopRecording()
				return
			}
		}

		time.Sleep(loopSleep)
	}
}

// Listen for a clap sounds to stop recording
func ListenForClapToStopRecording(a *AudioRecorder) {
	ListenForNClapsToStopRecording(a, 1)
}

// Listen for triple clap sounds to stop recording
func ListenForTripleClapsToStopRecording(a *AudioRecorder) {
	ListenForNClapsToStopRecording(a, 3)
}
