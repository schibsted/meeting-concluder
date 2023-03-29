package concluder

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type ChatGPT struct {
	apiKey string
}

func NewChatGPT(apiKey string) *ChatGPT {
	return &ChatGPT{
		apiKey: apiKey,
	}
}

func (cg *ChatGPT) GenerateConclusion(transcript string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	prompt := "Create a concise conclusion of the following meeting transcript:\n" + transcript

	data := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant that provides concise conclusions of meeting transcripts.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cg.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("ChatGPT API request failed with status: " + resp.Status)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return "", err
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", errors.New("No conclusion generated")
}
