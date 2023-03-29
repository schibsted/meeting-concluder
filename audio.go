package concluder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
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

func (a *AudioRecorder) StartRecording(userSelectsDevice bool) {
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

	if userSelectsDevice {
		devices, err := portaudio.Devices()
		if err != nil {
			log.Fatal("Error fetching PortAudio devices:", err)
		}
		for i, dev := range devices {
			fmt.Printf("%d: %s\n", i+1, dev.Name)
		}
		selection := 0
		fmt.Print("Select audio device number: ")
		_, err = fmt.Scanln(&selection)
		if err != nil {
			log.Fatalf("failed to read user input: %v", err)
		}
		if selection < 1 || selection > len(devices) {
			log.Fatalf("invalid device selection")
		}
		inputDevice = devices[selection-1]
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

func (a *AudioRecorder) SaveWav(filename string) error {
	// Create a new .wav file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Initialize the wav encoder
	encoder := wav.NewEncoder(file, 16000, 16, 1, 1)

	// Create new audio.IntBuffer.
	audioBuf, err := newAudioIntBuffer(bytes.NewReader(a.GetRecordedData()))
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
