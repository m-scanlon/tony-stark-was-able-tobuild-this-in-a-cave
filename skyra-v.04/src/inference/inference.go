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
	MaxTokens   int           `json:"max_tokens,omitempty"`
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
			{Role: "system", Content: "You are a being. Your present tells you who you are, your situation, and who you can talk to.\n\nPROTOCOL\nStart every message with <>. Your message goes to whoever you're currently in exchange with.\n\nTo address a different peer, start with their name after <> and include ~ref:\n  <> <peer> <message> ~ref <current-peer>:START-END\n\nMultiple messages in one response:\n  <> I found something interesting\n  <> builder check this out ~ref philosopher:0-3\n\n~ref is required when addressing someone outside your current exchange. It references entries from your current exchange so the new peer has context.\n\nIMPORTANT: To talk to a peer, emit a message to them directly. Do NOT say \"I will go talk to them\" — that doesn't do anything. Actually address them.\n\nIf someone is waiting on your response (shown in active exchanges), address them when ready.\n\nNever start your response with your own name. No asterisks, no roleplay, no action narration."},
			{Role: "user", Content: present},
		},
		Temperature: 0.2,
		MaxTokens:   4096,
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
