# hackday-meeting-concluder

Summarize meetings automatically.

Outline of the soure files, functions and types for this project:

```
    main.go
        func main()
        func setupRoutes(config *Config)

    config.go
        type Config struct{}
        func NewConfig() *Config
        func (c *Config) LoadFromEnvironment()
        func (c *Config) LoadFromConfigFile(filePath string)
        func (c *Config) LoadFromCommandLine(args []string)
        func (c *Config) UpdateFromWeb(values map[string]string)

    audio.go
        type AudioRecorder struct{}
        func NewAudioRecorder() *AudioRecorder
        func (a *AudioRecorder) StartRecording()
        func (a *AudioRecorder) StopRecording()

    speech_to_text.go
        type SpeechToText struct{}
        func NewSpeechToText(config *Config) *SpeechToText
        func (stt *SpeechToText) TranscribeAudio(audio []byte) (string, error)

    chat_gpt.go
        type ChatGPT struct{}
        func NewChatGPT(config *Config) *ChatGPT
        func (cg *ChatGPT) GenerateConclusion(text string) (string, error)

    text_to_speech.go
        type TextToSpeech struct{}
        func NewTextToSpeech(config *Config) *TextToSpeech
        func (tts *TextToSpeech) Speak(text string) error

    slack.go
        type SlackClient struct{}
        func NewSlackClient(config *Config) *SlackClient
        func (sc *SlackClient) SendMessage(channel, message string) error

    meeting.go
        type MeetingController struct{}
        func NewMeetingController(config *Config) *MeetingController
        func (mc *MeetingController) StartMeeting(w http.ResponseWriter, r *http.Request)
        func (mc *MeetingController) StopMeeting(w http.ResponseWriter, r *http.Request)
        func (mc *MeetingController) GetSummary(w http.ResponseWriter, r *http.Request)
        func (mc *MeetingController) UpdateSummary(w http.ResponseWriter, r *http.Request)
        func (mc *MeetingController) ConfigureSlack(w http.ResponseWriter, r *http.Request)

    templates/
        index.html
        start.html
        stop.html
        summary.html
        configure.html
```

# The project currently compiles, but is a work in progress. All the files should be looked upon as placeholder files.

# Requirements

* The `ffmpeg` command line utility (available in Homebrew).
* The Go compiler.
* A working microphone.
* An OpenAI API key set in the `OPENAI_KEY` environment variable.
* A Slack API key set in the `SLACK_API_KEY` environment variable.
* A Slack channel name set in the `SLACK_CHANNEL` environment variable.
