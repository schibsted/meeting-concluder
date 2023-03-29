package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	const inputFile = "input.wav"
	fmt.Printf("Playing %s...\n", inputFile)
	cmd := exec.Command("afplay", inputFile)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error playing audio:", err)
		os.Exit(1)
	}
}
