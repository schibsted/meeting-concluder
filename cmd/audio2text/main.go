// Copyright 2023 Schibsted. Licensed under the terms of the Apache 2.0 license. See LICENSE in the project root.

package main

import (
	"fmt"
	"os"

	concluder "github.com/schibsted/meeting-concluder"
)

const mp4FilePath = "input.mp4"

func main() {
	transcript, err := concluder.TranscribeAudio(mp4FilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error transribing %s: %v\n", mp4FilePath, err)
		return
	}

	transcript += "\n"
	fmt.Printf("%s\n", transcript)

	if err := os.WriteFile("output.txt", []byte(transcript), 0o644); err != nil {
		fmt.Printf("Error writing to output.txt: %v\n", err)
		return
	}
}
