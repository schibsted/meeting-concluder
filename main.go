package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Create and load the configuration
	config := NewConfig()
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

func setupRoutes(config *Config) *mux.Router {
	router := mux.NewRouter()

	// Create MeetingController
	meetingController := NewMeetingController(config)

	// Set up routes
	router.HandleFunc("/", meetingController.Index).Methods("GET")
	router.HandleFunc("/start", meetingController.StartMeeting).Methods("POST")
	router.HandleFunc("/stop", meetingController.StopMeeting).Methods("POST")
	router.HandleFunc("/summary", meetingController.GetSummary).Methods("GET")
	router.HandleFunc("/update-summary", meetingController.UpdateSummary).Methods("POST")
	router.HandleFunc("/configure", meetingController.ConfigureSlack).Methods("POST")

	return router
}
