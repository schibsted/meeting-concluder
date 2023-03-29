package concluder

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type SpeechToText struct {
	apiKey string
}

func NewSpeechToText(config *Config) *SpeechToText {
	return &SpeechToText{
		apiKey: config.WhisperAPIKey,
	}
}

func (stt *SpeechToText) TranscribeAudio(audio []byte) (string, error) {
	url := "https://api.openai.com/v1/whisper/asr"

	req, err := http.NewRequest("POST", url, bytes.NewReader(audio))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "audio/x-wav")
	req.Header.Set("Authorization", "Bearer "+stt.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Whisper API request failed with status: " + resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Text string `json:"text"`
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return "", err
	}

	return result.Text, nil
}
