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
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	part, err := writer.CreateFormFile("file", audioFilePath)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error copying file contents: %v", err)
	}
	err = writer.WriteField("model", "whisper-1")
	if err != nil {
		return "", fmt.Errorf("error writing form field: %v", err)
	}
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing multipart writer: %v", err)
	}

	req, err := http.NewRequest("POST", whisperURL, body)
	if err != nil {
		return "", fmt.Errorf("error creating API request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Config.OpenAIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending API request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading API response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: %s", resp.Status)
	}

	transcript, err := extractText(string(responseBody))
	if err != nil {
		return "", fmt.Errorf("error parsing the response: %v", err)
	}

	return transcript, nil
}
