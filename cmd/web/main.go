package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xyproto/env/v2"
	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

type RecordRequest struct {
	SlackChannel string `json:"slack_channel"`
	SlackAPIKey  string `json:"slack_api_key"`
	OpenAIAPIKey string `json:"openai_api_key"`
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	router.Post("/record", func(w http.ResponseWriter, r *http.Request) {
		var req RecordRequest
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &req)

		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		summary, err := concluder.RecordAndConclude(req.SlackChannel, req.SlackAPIKey, req.OpenAIAPIKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"summary": summary})
	})

	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	addr := env.Str("HOST", ":3000")

	log.Println("Starting server on " + addr)
	http.ListenAndServe(addr, router)
}
