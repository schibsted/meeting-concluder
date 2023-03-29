package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

func main() {
	// Initialize the AudioRecorder from the concluder package
	audioRecorder := concluder.NewAudioRecorder()

	// Start recording
	audioRecorder.StartRecording()

	// Record for 5 seconds
	fmt.Println("Recording for 5 seconds...")
	time.Sleep(5 * time.Second)

	// Stop recording
	audioRecorder.StopRecording()

	// Get the recorded data
	recordedData := audioRecorder.GetRecordedData()

	// Save the recorded data to a .wav file
	err := saveToWavFile("output.wav", recordedData)
	if err != nil {
		fmt.Println("Error saving .wav file:", err)
		os.Exit(1)
	}

	fmt.Println("Audio saved to output.wav")
}

func saveToWavFile(filename string, data []byte) error {
	// Create a new .wav file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Initialize the wav encoder
	encoder := wav.NewEncoder(file, 16000, 16, 1, 1)

	// Create new audio.IntBuffer.
	audioBuf, err := newAudioIntBuffer(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	// Write the recorded data to the .wav file
	err = encoder.Write(audioBuf)
	if err != nil {
		return err
	}

	// Close the .wav file
	err = encoder.Close()
	if err != nil {
		return err
	}

	return nil
}

func newAudioIntBuffer(r io.Reader) (*audio.IntBuffer, error) {
	buf := audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  8000,
		},
	}
	for {
		var sample int16
		err := binary.Read(r, binary.LittleEndian, &sample)
		switch {
		case err == io.EOF:
			return &buf, nil
		case err != nil:
			return nil, err
		}
		buf.Data = append(buf.Data, int(sample))
	}
}
