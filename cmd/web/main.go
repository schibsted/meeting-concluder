package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

func main() {
	// Create and load the configuration
	config := concluder.NewConfig()
	config.LoadFromEnvironment()
	config.LoadFromConfigFile("config.json")
	config.LoadFromCommandLine(os.Args)

	// Set up routes
	router := setupRoutes(config)

	// Start the web server
	log.Println("Starting the web server on :8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func setupRoutes(config *concluder.Config) http.Handler {
	router := chi.NewRouter()

	// Create MeetingController
	meetingController := concluder.NewMeetingController(config)

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set up routes
	router.Get("/", meetingController.Index)
	router.Post("/start", meetingController.StartMeeting)
	router.Post("/stop", meetingController.StopMeeting)
	router.Get("/summary", meetingController.GetSummary)
	router.Post("/update-summary", meetingController.UpdateSummary)
	router.Post("/configure", meetingController.ConfigureSlack)

	return router
}
