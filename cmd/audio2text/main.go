package main

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

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

const (
	whisperURL  = "https://api.openai.com/v1/audio/transcriptions"
	mp4FilePath = "input.mp4"
)

type AsrResponse struct {
	Transcript string `json:"transcript"`
}

type Response struct {
	Text string `json:"text"`
}

func extractText(jsonStr string) (string, error) {
	var response Response
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return "", err
	}
	return response.Text, nil
}

func main() {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open(mp4FilePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()
	part, err := writer.CreateFormFile("file", mp4FilePath)
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file contents: %v\n", err)
		return
	}
	err = writer.WriteField("model", "whisper-1")
	if err != nil {
		fmt.Printf("Error writing form field: %v\n", err)
		return
	}
	err = writer.Close()
	if err != nil {
		fmt.Printf("Error closing multipart writer: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", whisperURL, body)
	if err != nil {
		fmt.Printf("Error creating API request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", concluder.Config.OpenAI_APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending API request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading API response: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", resp.Status)
		return
	}

	responseString := string(responseBody)
	fmt.Printf("Response: %s\n", responseString)

	transcript, err := extractText(responseString)
	if err != nil {
		fmt.Printf("Error parsing the response: %v\n", err)
		return
	}

	fmt.Printf("Transcript: %s\n", transcript)

	if !strings.HasSuffix(transcript, "\n") {
		transcript += "\n"
	}

	if err := os.WriteFile("output.txt", []byte(transcript), 0o644); err != nil {
		fmt.Printf("Error writing to output.txt: %v\n", err)
		return
	}
}
