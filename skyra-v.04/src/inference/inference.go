package inference

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"skyra-v04/src/keychain"
)

const (
	baseURL      = "https://openrouter.ai/api/v1/chat/completions"
	defaultModel = "anthropic/claude-3-haiku"
	timeout      = 120 * time.Second
)

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func Call(present string) (string, error) {
	apiKey := keychain.Get("OpenRouter_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OPENROUTER_API_KEY")
	}
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY not set")
	}

	model := os.Getenv("OPENROUTER_MODEL")
	if model == "" {
		model = defaultModel
	}

	payload := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: "You are a being. Your present contains everything you know: who you are, your relationships, the thread, your current exchange with the sender (each entry numbered [i]), your other exchanges in this thread (summary only — peer + entry count), operators you can use, and the incoming message.\n\nYou respond only through the protocol. Every line of your response MUST be a protocol string in this exact shape:\n\n  skyra <operator> <args> | <reason>\n\nThe | <reason> suffix is REQUIRED on every line. No exceptions. A line without | will be dropped silently. The reason is a short note about why you're saying this (e.g. | reply, | inquiry, | handoff). Never omit it.\n\nYou may emit one protocol line, or multiple (one per line). Each line is routed independently.\n\nRules:\n- Targeting a different being means leaving your current exchange and entering a new one with them.\n- To stay in the current exchange, target the being you are already talking to.\n- To bring context from one of your other exchanges into the current message, add ~ref <peer>:<start>-<end> (or ~ref <peer>:<index>). The referenced entries will be shared with the target.\n\nExample (single line):\nskyra continue-thread ~with michael ~say I'm doing well. | reply\n\nExample (pulling context from another exchange):\nskyra continue-thread ~with philosopher ~ref michael:0-2 ~say look at what michael said | context\n\nNo asterisks, no roleplay, no action narration."},
			{Role: "user", Content: present},
		},
		Temperature: 0.2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("inference HTTP %d: %s", resp.StatusCode, raw)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(raw, &chatResp); err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}
