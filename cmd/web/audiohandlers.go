package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

type AudioHandlers struct {
	audioRecorder *concluder.AudioRecorder
}

func NewAudioHandlers(ar *concluder.AudioRecorder) *AudioHandlers {
	return &AudioHandlers{audioRecorder: ar}
}

func (ah *AudioHandlers) startRecordingHandler(w http.ResponseWriter, r *http.Request) {
	if ah.audioRecorder.IsRecording() {
		http.Error(w, "Already recording", http.StatusConflict)
		return
	}
	ah.audioRecorder.StartRecording()
	w.WriteHeader(http.StatusOK)
}

func (ah *AudioHandlers) stopRecordingHandler(w http.ResponseWriter, r *http.Request) {
	if !ah.audioRecorder.IsRecording() {
		http.Error(w, "Not recording", http.StatusConflict)
		return
	}
	ah.audioRecorder.StopRecording()
	filename := fmt.Sprintf("output_%s.wav", time.Now().Format("2006-01-02_15-04-05"))
	if err := ah.audioRecorder.SaveWav(filename); err != nil {
		http.Error(w, "Error saving recording", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"filename": filename})
}
