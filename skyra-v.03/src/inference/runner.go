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

	"skyra-v03/src/metaxu"
)

const (
	maxRetries     = 5
	baseRetryDelay = 2 * time.Second
)

const systemPrompt = "You are a being in a cognitive system. Always respond with a single protocol string in this exact format: skyra <being> <expression> | <source>: <reason>. No explanation. No markdown. No extra text."

type Config struct {
	BaseURL string
	Model   string
}

type Runner struct {
	config Config
	client *http.Client
}

func New(config Config) *Runner {
	return &Runner{
		config: config,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	Stream      bool          `json:"stream"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

func (r *Runner) Run(present string, originName string) (metaxu.Signal, error) {
	fmt.Fprintf(os.Stderr, "[inference] being=%s\n", originName)
	fmt.Fprintf(os.Stderr, "[present]\n%s\n[/present]\n", present)
	protocol, err := r.callWithRetry(present)
	if err != nil {
		return metaxu.Signal{}, fmt.Errorf("inference: %w", err)
	}
	fmt.Fprintf(os.Stderr, "[inference] response=%s\n", protocol)

	return metaxu.Signal{
		Origin:  originName,
		Impulse: protocol,
	}, nil
}

func (r *Runner) callWithRetry(present string) (string, error) {
	delay := baseRetryDelay
	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, err := r.call(present)
		if err == nil {
			return result, nil
		}
		if !strings.HasPrefix(err.Error(), "retry:") {
			return "", err
		}
		if attempt == maxRetries {
			return "", fmt.Errorf("inference unavailable after %d retries: %w", maxRetries, err)
		}
		fmt.Fprintf(os.Stderr, "[inference] retrying in %s (attempt %d/%d)\n", delay, attempt+1, maxRetries)
		time.Sleep(delay)
		delay *= 2
	}
	return "", fmt.Errorf("unreachable")
}

func (r *Runner) call(present string) (string, error) {
	endpoint := r.config.BaseURL + "/v1/chat/completions"

	payload := chatRequest{
		Model: r.config.Model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: present},
		},
		Temperature: 0.2,
		Stream:      false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == 429 {
		return "", fmt.Errorf("retry:%d:%s", resp.StatusCode, rawBody)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("inference HTTP %d: %s", resp.StatusCode, rawBody)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(rawBody, &chatResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if len(chatResp.Choices) == 0 || strings.TrimSpace(chatResp.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("empty response from model")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}
