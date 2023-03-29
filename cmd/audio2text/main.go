package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
)

const (
	whisperURL  = "https://api.openai.com/v1/audio/transcriptions"
	mp4FilePath = "output.mp4"
)

type AsrResponse struct {
	Transcript string `json:"transcript"`
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

	fmt.Printf("Transcript: %s\n", string(responseBody))
}
