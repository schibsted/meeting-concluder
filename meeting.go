package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type MeetingController struct {
	audioRecorder  *AudioRecorder
	speechToText   *SpeechToText
	chatGPT        *ChatGPT
	textToSpeech   *TextToSpeech
	slackClient    *SlackClient
	meetingStarted bool
	startTime      time.Time
}

func NewMeetingController(config *Config) *MeetingController {
	return &MeetingController{
		audioRecorder: NewAudioRecorder(),
		speechToText:  NewSpeechToText(config),
		chatGPT:       NewChatGPT(config),
		textToSpeech:  NewTextToSpeech(config),
		slackClient:   NewSlackClient(config),
	}
}

func (mc *MeetingController) StartMeeting(w http.ResponseWriter, r *http.Request) {
	if !mc.meetingStarted {
		mc.meetingStarted = true
		mc.startTime = time.Now()
		mc.audioRecorder.StartRecording()
	}
}

func (mc *MeetingController) StopMeeting(w http.ResponseWriter, r *http.Request) {
	if mc.meetingStarted {
		audioData := mc.audioRecorder.StopRecording()
		mc.meetingStarted = false
		transcription, err := mc.speechToText.TranscribeAudio(audioData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		conclusion, err := mc.chatGPT.GenerateConclusion(transcription)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mc.textToSpeech.Speak(conclusion)
	}
}

func (mc *MeetingController) GetSummary(w http.ResponseWriter, r *http.Request) {
	audioData := mc.audioRecorder.GetRecordedData()
	transcription, err := mc.speechToText.TranscribeAudio(audioData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	conclusion, err := mc.chatGPT.GenerateConclusion(transcription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"conclusion": conclusion})
}

func (mc *MeetingController) UpdateSummary(w http.ResponseWriter, r *http.Request) {
	updatedConclusion := r.FormValue("conclusion")
	mc.textToSpeech.Speak(updatedConclusion)
	mc.slackClient.SendMessage("#nmp-meeting-concluder", updatedConclusion)
}

func (mc *MeetingController) ConfigureSlack(w http.ResponseWriter, r *http.Request) {
	slackToken := r.FormValue("token")
	channel := r.FormValue("channel")
	mc.slackClient.UpdateConfig(slackToken, channel)
}

func (mc *MeetingController) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/start", http.StatusSeeOther)
}

func (mc *MeetingController) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", mc.Index)
	router.Post("/start", mc.StartMeeting)
	router.Post("/stop", mc.StopMeeting)
	router.Get("/summary", mc.GetSummary)
	router.Post("/update-summary", mc.UpdateSummary)
	router.Post("/configure", mc.ConfigureSlack)

	return router
}
