package concluder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const whisperURL = "https://api.openai.com/v1/audio/transcriptions"

type AsrResponse struct {
	Transcript string `json:"transcript"`
}

type AudioResponse struct {
	Text string `json:"text"`
}

func extractText(jsonStr string) (string, error) {
	var response AudioResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return "", err
	}
	return response.Text, nil
}

// TranscribeAudio can extract the text from .mp4, .wav or several different audio file types
func TranscribeAudio(audioFilePath string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open(audioFilePath)
	if err != nil {
		return "", fmt.Errorf("Error opening file: %v\n", err)
	}
	defer file.Close()
	part, err := writer.CreateFormFile("file", audioFilePath)
	if err != nil {
		return "", fmt.Errorf("Error creating form file: %v\n", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("Error copying file contents: %v\n", err)
	}
	err = writer.WriteField("model", "whisper-1")
	if err != nil {
		return "", fmt.Errorf("Error writing form field: %v\n", err)
	}
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("Error closing multipart writer: %v\n", err)
	}

	req, err := http.NewRequest("POST", whisperURL, body)
	if err != nil {
		return "", fmt.Errorf("Error creating API request: %v\n", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Config.OpenAI_APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending API request: %v\n", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading API response: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error: %s\n", resp.Status)
	}

	transcript, err := extractText(string(responseBody))
	if err != nil {
		return "", fmt.Errorf("Error parsing the response: %v\n", err)
	}

	if !strings.HasSuffix(transcript, "\n") {
		transcript += "\n"
	}

	return transcript, nil
}
