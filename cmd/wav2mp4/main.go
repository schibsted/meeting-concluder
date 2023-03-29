package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	inputFile := "input.wav"
	outputFile := "output.mp4"

	err := convertToMP4(inputFile, outputFile)
	if err != nil {
		log.Fatalf("Error converting %s to %s: %v", inputFile, outputFile, err)
	}

	log.Printf("Successfully converted %s to %s", inputFile, outputFile)
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
