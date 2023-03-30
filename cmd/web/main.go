package main

import (
	"fmt"
	"log"
	"net/http"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
	"github.com/xyproto/env/v2"
)

func main() {
	ar := concluder.NewAudioRecorder()
	defer ar.Done()

	err := ar.UserSelectsTheInputDevice()
	if err != nil {
		log.Fatal(err)
	}

	audioHandlers := NewAudioHandlers(ar)

	http.HandleFunc("/start", audioHandlers.startRecordingHandler)
	http.HandleFunc("/stop", audioHandlers.stopRecordingHandler)

	addr := env.Str("HOST", ":3000")
	fmt.Printf("Starting server on %s...\n", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
