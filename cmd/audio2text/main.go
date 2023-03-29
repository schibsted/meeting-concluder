package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/xyproto/env/v2"
)

const (
	apiKEy      = env.StrAlt("OPENAI_API_KEY", "OPENAI_KEY", env.Str("CHATGPT_API_KEY"))
	whisperURL  = "https://api.openai.com/v1/whisper/asr"
	wavFilePath = "output.wav"
)

type AsrParams struct {
	Engine string `url:"engine"`
}

type AsrResponse struct {
	Transcript string `json:"transcript"`
}

func main() {
	// Read the WAV file
	wavData, err := ioutil.ReadFile(wavFilePath)
	if err != nil {
		fmt.Printf("Error reading WAV file: %v\n", err)
		return
	}

	// Prepare the API request
	client := &http.Client{}
	req, err := http.NewRequest("POST", whisperURL, bytes.NewReader(wavData))
	if err != nil {
		fmt.Printf("Error creating API request: %v\n", err)
		return
	}

	// Set the required headers
	req.Header.Set("Content-Type", "audio/wav")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Add the query parameters
	params := AsrParams{Engine: "whisper-asr"}
	q, _ := query.Values(params)
	req.URL.RawQuery = q.Encode()

	// Send the API request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending API request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", resp.Status)
		return
	}

	// Parse the API response
	var asrResponse AsrResponse
	err = json.NewDecoder(resp.Body).Decode(&asrResponse)
	if err != nil {
		fmt.Printf("Error decoding API response: %v\n", err)
		return
	}

	// Print the transcript
	fmt.Printf("Transcript: %s\n", asrResponse.Transcript)
}
