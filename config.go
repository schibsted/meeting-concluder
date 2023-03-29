package concluder

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

type Config struct {
	WhisperAPIKey      string
	ChatGPTAPIKey      string
	TextToSpeechAPIKey string
	SlackToken         string
	SlackChannel       string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadFromEnvironment() {
	c.WhisperAPIKey = os.Getenv("WHISPER_API_KEY")
	c.ChatGPTAPIKey = os.Getenv("CHAT_GPT_API_KEY")
	c.TextToSpeechAPIKey = os.Getenv("TEXT_TO_SPEECH_API_KEY")
	c.SlackToken = os.Getenv("SLACK_TOKEN")
	c.SlackChannel = os.Getenv("SLACK_CHANNEL")
}

func (c *Config) LoadFromConfigFile(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return
	}
}

func (c *Config) LoadFromCommandLine(args []string) {
	flagSet := flag.NewFlagSet("config", flag.ContinueOnError)
	flagSet.StringVar(&c.WhisperAPIKey, "whisper-api-key", c.WhisperAPIKey, "Whisper API Key")
	flagSet.StringVar(&c.ChatGPTAPIKey, "chat-gpt-api-key", c.ChatGPTAPIKey, "ChatGPT API Key")
	flagSet.StringVar(&c.TextToSpeechAPIKey, "text-to-speech-api-key", c.TextToSpeechAPIKey, "Text to Speech API Key")
	flagSet.StringVar(&c.SlackToken, "slack-token", c.SlackToken, "Slack Token")
	flagSet.StringVar(&c.SlackChannel, "slack-channel", c.SlackChannel, "Slack Channel")
	flagSet.Parse(args[1:])
}

func (c *Config) UpdateFromWeb(values map[string]string) {
	if apiKey, ok := values["whisper_api_key"]; ok {
		c.WhisperAPIKey = apiKey
	}
	if apiKey, ok := values["chat_gpt_api_key"]; ok {
		c.ChatGPTAPIKey = apiKey
	}
	if apiKey, ok := values["text_to_speech_api_key"]; ok {
		c.TextToSpeechAPIKey = apiKey
	}
	if token, ok := values["slack_token"]; ok {
		c.SlackToken = token
	}
	if channel, ok := values["slack_channel"]; ok {
		c.SlackChannel = channel
	}
}
