package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type OpenAIClient struct {
	APIKey string
	URL    string
	Model  string
	Client *http.Client
}

func NewOpenAIClient(apiKey, url, model string) *OpenAIClient {
	return &OpenAIClient{
		APIKey: apiKey,
		URL:    url,
		Model:  model,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (o *OpenAIClient) Generate(
	ctx context.Context,
	systemPrompt string,
	input map[string]interface{},
	temperature float64,
) (map[string]interface{}, error) {

	userInputBytes, _ := json.Marshal(input)

	reqBody := chatRequest{
		Model: o.Model,
		Messages: []chatMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: string(userInputBytes),
			},
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		o.URL,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+o.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("LLM API error")
	}

	var parsed chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	if len(parsed.Choices) == 0 {
		return nil, errors.New("empty LLM response")
	}

	content := parsed.Choices[0].Message.Content

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, errors.New("LLM returned invalid JSON")
	}

	return result, nil
}
