package concluder

import (
	"bytes"
	"log"

	"github.com/gordonklaus/portaudio"
)

type AudioRecorder struct {
	stream    *portaudio.Stream
	buffer    bytes.Buffer
	recording bool
}

func NewAudioRecorder() *AudioRecorder {
	return &AudioRecorder{}
}

func (a *AudioRecorder) StartRecording() {
	if a.recording {
		log.Println("Audio recording is already in progress.")
		return
	}

	err := portaudio.Initialize()
	if err != nil {
		log.Fatal("Error initializing PortAudio:", err)
	}

	inputDevice, err := portaudio.DefaultInputDevice()
	if err != nil {
		log.Fatal("Error fetching default input device:", err)
	}

	parameters := portaudio.LowLatencyParameters(inputDevice, nil)
	parameters.Input.Channels = 1
	parameters.SampleRate = 16000
	parameters.FramesPerBuffer = 1024

	stream, err := portaudio.OpenStream(parameters, a.captureAudio)
	if err != nil {
		log.Fatal("Error opening PortAudio stream:", err)
	}

	err = stream.Start()
	if err != nil {
		log.Fatal("Error starting PortAudio stream:", err)
	}

	a.stream = stream
	a.recording = true
	log.Println("Started audio recording.")
}

func (a *AudioRecorder) StopRecording() {
	if !a.recording {
		return
	}

	a.recording = false
	a.stream.Stop()
	a.stream.Close()
	portaudio.Terminate()

	log.Println("Stopped audio recording.")
}

func (a *AudioRecorder) GetRecordedData() []byte {
	return a.buffer.Bytes()
}

func (a *AudioRecorder) captureAudio(inputBuffer, _ []float32) {
	if !a.recording {
		return
	}

	for _, sample := range inputBuffer {
		encodedSample := int16(sample * 32767)
		a.buffer.WriteByte(byte(encodedSample))
		a.buffer.WriteByte(byte(encodedSample >> 8))
	}
}
