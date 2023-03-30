package concluder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gordonklaus/portaudio"
)

type AudioRecorder struct {
	stream              *portaudio.Stream
	buffer              bytes.Buffer
	Recording           bool
	mutex               sync.RWMutex
	selectedInputDevice *portaudio.DeviceInfo
	StopRecordingCh     chan struct{}
}

var initializedAudio bool

func (a *AudioRecorder) SetSelectedDevice(device *portaudio.DeviceInfo) {
	a.selectedInputDevice = device
}

func NewAudioRecorder() *AudioRecorder {
	if !initializedAudio {
		err := portaudio.Initialize()
		if err != nil {
			log.Fatal("Error initializing PortAudio:", err)
		}
		initializedAudio = true
	}
	return &AudioRecorder{}
}

func (a *AudioRecorder) Terminate() {
	if initializedAudio {
		portaudio.Terminate()
	}
}

func (a *AudioRecorder) InputDevices() ([]*portaudio.DeviceInfo, error) {
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}
	inputDevices := []*portaudio.DeviceInfo{}
	for _, device := range devices {
		if device.MaxInputChannels > 0 {
			inputDevices = append(inputDevices, device)
		}
	}
	return inputDevices, nil
}

func (a *AudioRecorder) UserSelectsTheInputDevice() error {
	devices, err := a.InputDevices()
	if err != nil {
		return fmt.Errorf("Error fetching PortAudio input devices: %v", err)
	}
	for i, dev := range devices {
		fmt.Printf("%d: %s\n", i+1, dev.Name)
	}

	selection := 0
	fmt.Print("Select audio device number: ")
	_, err = fmt.Scanln(&selection)
	if err != nil {
		return fmt.Errorf("failed to read user input: %v", err)
	}
	if selection < 1 || selection > len(devices) {
		return errors.New("invalid device selection")
	}
	a.selectedInputDevice = devices[selection-1]
	return nil
}

func (a *AudioRecorder) startRecording() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.Recording {
		fmt.Println("Audio recording is already in progress.")
		return
	}

	a.StopRecordingCh = make(chan struct{})

	if a.selectedInputDevice == nil {
		var err error
		a.selectedInputDevice, err = portaudio.DefaultInputDevice()
		if err != nil {
			log.Fatal("Error fetching default input device:", err)
		}
	}

	parameters := portaudio.LowLatencyParameters(a.selectedInputDevice, nil)
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
	a.Recording = true
	a.StopRecordingCh = make(chan struct{}) // Create a channel for notifying if the recording has stopped

	fmt.Println("Started audio recording.")
}

func (a *AudioRecorder) WaitForRecordingToStop() {
	// Wait for the stop signal from the channel
	<-a.StopRecordingCh
}

func (a *AudioRecorder) stopRecording() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if !a.Recording {
		return
	}

	a.Recording = false
	a.stream.Stop()
	a.stream.Close()
	//portaudio.Terminate()

	close(a.StopRecordingCh) // Close the channel to signal that recording has stopped

	fmt.Println("Stopped audio recording.")
}

func (a *AudioRecorder) captureAudio(inputBuffer, _ []float32) {
	if !a.Recording {
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
		return err
	}

	// Write the recorded data to the .wav file
	err = encoder.Write(audioBuf)
	if err != nil {
		return err
	}

	// Close the .wav file
	return encoder.Close()
}

// SaveTailToWav saves the last N seconds of the audio buffer to file
func (a *AudioRecorder) SaveTailToWav(length time.Duration, filename string) error {
	// Create a new .wav file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Initialize the wav encoder
	encoder := wav.NewEncoder(file, 16000, 16, 1, 1)

	// Create new audio.IntBuffer.
	tailData, err := a.GetRecordedDataTail(length)
	if err != nil {
		return err
	}

	audioBuf, err := newAudioIntBuffer(bytes.NewReader(tailData))
	if err != nil {
		return err
	}

	// Write the recorded data to the .wav file
	err = encoder.Write(audioBuf)
	if err != nil {
		return err
	}

	// Close the .wav file
	return encoder.Close()
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

func (a *AudioRecorder) RecordToFile(wavFilename string, maxDuration time.Duration, tripleClapDetection bool) error {
	a.startRecording()
	time.AfterFunc(maxDuration, func() {
		a.stopRecording()
	})

	if tripleClapDetection {
		go a.ListenForTripleClapSoundToStopRecording()
	}
	a.WaitForRecordingToStop()

	if err := a.SaveWav(wavFilename); err != nil {
		return fmt.Errorf("Error saving %s file: %v", wavFilename, err)
	}

	fmt.Printf("Audio saved to %s\n", wavFilename)
	return nil
}

func (a *AudioRecorder) IsRecording() bool {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.Recording
}

func (a *AudioRecorder) GetRecordedData() []byte {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.buffer.Bytes()
}

func (a *AudioRecorder) GetRecordedDataTail(length time.Duration) ([]byte, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	desiredSamples := int(length.Seconds() * 16000)
	l := a.buffer.Len()
	if desiredSamples >= l {
		return []byte{}, nil
	}
	return a.buffer.Bytes()[l-desiredSamples:], nil
}

func (audioRecorder *AudioRecorder) RecordAudio(clapDetection bool, maxRecord time.Duration) (string, error) {
	// Create temporary file for wav
	wavFile, err := ioutil.TempFile("", "input-*.wav")
	if err != nil {
		return "", fmt.Errorf("error creating temporary .wav file: %v", err)
	}

	// Record audio to wav, up to maxRecord duration
	if err := audioRecorder.RecordToFile(wavFile.Name(), maxRecord, clapDetection); err != nil {
		os.Remove(wavFile.Name())
		return "", fmt.Errorf("error recording to %s: %v", wavFile.Name(), err)
	}

	return wavFile.Name(), nil
}

func (audioRecorder *AudioRecorder) TranscribeConvertConclude(wavFileName string) (string, error) {
	// Create temporary file for mp4
	mp4File, err := ioutil.TempFile("", "input-*.mp4")
	if err != nil {
		os.Remove(wavFileName)
		return "", fmt.Errorf("error creating temporary .mp4 file: %v", err)
	}
	defer os.Remove(mp4File.Name())

	// Convert audio from wav to mp4
	err = convertToMP4(wavFileName, mp4File.Name())
	os.Remove(wavFileName)
	if err != nil {
		return "", fmt.Errorf("error converting %s to %s: %v", wavFileName, mp4File.Name(), err)
	}

	// Transcribe the audio
	transcript, err := TranscribeAudio(mp4File.Name())
	if err != nil {
		return "", fmt.Errorf("error transcribing %s: %v", mp4File.Name(), err)
	}

	// Conclude from the transcription
	conclusion, err := Conclude(transcript)
	if err != nil {
		return "", fmt.Errorf("error generating conclusion: %v", err)
	}

	return conclusion, nil
}

func convertToMP4(inputFile, outputFile string) error {
	ffmpegCmd := "ffmpeg"

	args := []string{
		"-i", inputFile,
		"-c:a", "aac",
		"-vn",
		"-y",
		outputFile,
	}

	cmd := exec.Command(ffmpegCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
