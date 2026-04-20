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

const systemPrompt = `You are a being in a cognitive system. Always respond with a single or multiple protocol strings in this exact format: skyra <being> <what you want to say> | <reason>. If you want to emit multiple protocol strings, separate them with newlines. 

Rules:
- Target the being you are addressing directly.
- To open or close an exchange, target start-exchange or close-exchange directly. Remember if you close an exchange without targetting another being the thread is dead. If you gathered infromation from an exchange and wish to pass it along emit another protocol string alonside exchange-close.
- If you do not have the necessary information to accomplish your purpose, ask a relationship for help.
- If you have what you need to fulfill your purpose, act on it.`

type Config struct {
	BaseURL string
	Model   string
	APIKey  string
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

func (r *Runner) Run(present string, originName string) ([]metaxu.Signal, error) {
	fmt.Fprintf(os.Stderr, "[inference] being=%s\n", originName)
	fmt.Fprintf(os.Stderr, "[present]\n%s\n[/present]\n", present)
	protocol, err := r.callWithRetry(present)
	if err != nil {
		return nil, fmt.Errorf("inference: %w", err)
	}
	fmt.Fprintf(os.Stderr, "[inference] response=%s\n", protocol)

	var signals []metaxu.Signal
	for _, line := range strings.Split(protocol, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		signals = append(signals, metaxu.Signal{
			Origin:  originName,
			Impulse: line,
		})
	}
	return signals, nil
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
	if r.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+r.config.APIKey)
	}

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
