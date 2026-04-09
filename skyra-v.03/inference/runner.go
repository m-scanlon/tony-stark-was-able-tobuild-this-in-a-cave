package inference

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"skyra-v03/src/domain"
)

const (
	geminiEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent"
	systemPrompt   = "You are this being."
)

type Config struct {
	APIKey string
	Model  string
}

type Runner struct {
	config Config
	client *http.Client
}

func New(config Config) *Runner {
	return &Runner{
		config: config,
		client: &http.Client{Timeout: 45 * time.Second},
	}
}

type geminiRequest struct {
	SystemInstruction systemInstruction `json:"system_instruction"`
	Contents          []content         `json:"contents"`
	GenerationConfig  generationConfig  `json:"generationConfig"`
}

type systemInstruction struct {
	Parts []part `json:"parts"`
}

type content struct {
	Role  string `json:"role"`
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

type generationConfig struct {
	Temperature float64 `json:"temperature"`
}

type geminiResponse struct {
	Candidates []candidate `json:"candidates"`
}

type candidate struct {
	Content content `json:"content"`
}

func (r *Runner) Run(present string, originName string) (domain.Signal, error) {
	protocol, err := r.callGemini(present)
	if err != nil {
		return domain.Signal{}, fmt.Errorf("inference: %w", err)
	}

	return domain.Signal{
		Origin:  originName,
		Impulse: protocol,
	}, nil
}

func (r *Runner) callGemini(present string) (string, error) {
	endpoint := fmt.Sprintf(geminiEndpoint, url.PathEscape(r.config.Model))

	payload := geminiRequest{
		SystemInstruction: systemInstruction{
			Parts: []part{{Text: systemPrompt}},
		},
		Contents: []content{
			{
				Role:  "user",
				Parts: []part{{Text: present}},
			},
		},
		GenerationConfig: generationConfig{Temperature: 0.2},
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
	req.Header.Set("x-goog-api-key", r.config.APIKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini HTTP %d: %s", resp.StatusCode, rawBody)
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(rawBody, &geminiResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	protocol := extractText(geminiResp)
	if protocol == "" {
		return "", fmt.Errorf("empty response from model")
	}

	return strings.TrimSpace(protocol), nil
}

func extractText(resp geminiResponse) string {
	if len(resp.Candidates) == 0 {
		return ""
	}
	if len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	return resp.Candidates[0].Content.Parts[0].Text
}
