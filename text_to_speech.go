package concluder

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type TextToSpeech struct {
	apiKey string
}

func NewTextToSpeech(config *Config) *TextToSpeech {
	return &TextToSpeech{
		apiKey: config.TextToSpeechAPIKey,
	}
}

func (tts *TextToSpeech) Speak(text string) error {
	audioData, err := tts.generateAudio(text)
	if err != nil {
		return err
	}

	tempFileName := fmt.Sprintf("temp_audio_%d.wav", time.Now().UnixNano())
	err = ioutil.WriteFile(tempFileName, audioData, 0644)
	if err != nil {
		return err
	}

	defer func() {
		_ = os.Remove(tempFileName)
	}()

	cmd := exec.Command("afplay", tempFileName)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (tts *TextToSpeech) generateAudio(text string) ([]byte, error) {
	url := "https://text-to-speech-api.example.com/v1/synthesize"

	data := map[string]string{
		"text": text,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tts.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Text-to-Speech API request failed with status: " + resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
