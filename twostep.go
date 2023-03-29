package concluder

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

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

func (audioRecorder *AudioRecorder) TranscribeConvertConcludePost(wavFileName string) (string, error) {
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

	// Post the conclusion to Slack
	err = SendMessage(conclusion)
	if err != nil {
		return "", fmt.Errorf("error sending message to Slack channel: %v", err)
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
