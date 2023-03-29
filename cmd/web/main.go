package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xyproto/env/v2"
	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

type RecordRequest struct {
	ClapDetection bool          `json:"clap_detection"`
	MaxDuration   time.Duration `json:"-"`
}

var stopRecordingChan chan struct{}
var audioRecorder = concluder.NewAudioRecorder()

func startRecordingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST method required", http.StatusMethodNotAllowed)
		return
	}

	// Create a new channel for signaling stop recording
	stopRecordingChan = make(chan struct{})

	clapDetection, _ := strconv.ParseBool(r.FormValue("clapDetection"))
	maxDuration, _ := time.ParseDuration(r.FormValue("maxDuration"))

	go func() {
		// Record audio with given settings and wait for stop signal
		wavFileName, err := audioRecorder.RecordAudio(clapDetection, maxDuration)
		if err != nil {
			log.Printf("Error recording audio: %v", err)
			return
		}
		select {
		case <-stopRecordingChan:
			log.Println("Stopped audio recording.")
		case <-time.After(maxDuration):
			log.Println("Max recording duration reached.")
		}
		if err := os.Remove(wavFileName); err != nil {
			log.Printf("Error removing temporary .wav file: %v", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func stopRecordingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST method required", http.StatusMethodNotAllowed)
		return
	}

	if stopRecordingChan != nil {
		close(stopRecordingChan)
		stopRecordingChan = nil
	} else {
		http.Error(w, "No recording to stop", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		indexTemplate := template.Must(template.ParseFiles("static/index.html"))
		clapDetection := env.Bool("CLAP_DETECTION")
		data := struct {
			ClapDetection bool
			MaxDuration   time.Duration
		}{clapDetection, 1 * time.Hour}
		indexTemplate.Execute(w, data)
	})

	router.Get("/main.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/main.js")
	})

	router.Get("/tailwind-3.3.0.min.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/tailwind-3.3.0.min.js")
	})

	router.Post("/record/start", startRecordingHandler)
	router.Post("/record/stop", stopRecordingHandler)

	addr := env.Str("HOST", ":3000")

	log.Println("Starting server on " + addr)
	http.ListenAndServe(addr, router)
}
