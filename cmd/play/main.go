package main

import (
	"fmt"
	"os"
	"os/exec"
)

const inputFile = "output.wav"

func main() {
	// Check if input file exists
	_, err := os.Stat(inputFile)
	if os.IsNotExist(err) {
		fmt.Printf("Error: %s does not exist\n", inputFile)
		os.Exit(1)
	}

	// Play audio file
	fmt.Printf("Playing %s...\n", inputFile)
	cmd := exec.Command("afplay", inputFile)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error playing audio:", err)
		os.Exit(1)
	}
}
