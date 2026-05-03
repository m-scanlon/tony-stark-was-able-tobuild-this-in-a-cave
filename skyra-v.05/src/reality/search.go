package reality

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Search struct {
	id string
}

func (s *Search) ID() string { return s.id }

func (s *Search) Create(r *Relation) Reality {
	return &Search{id: "search"}
}

func (s *Search) Realize(r *Relation) string {
	query := strings.TrimSpace(r.Impulse)
	if query == "" {
		return "no search query provided"
	}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey == "" {
		if r.Log != nil {
			r.Log("[search]: no FIRECRAWL_API_KEY")
		}
		return "search unavailable"
	}

	body, _ := json.Marshal(map[string]any{
		"query": query,
		"limit": 5,
	})

	req, err := http.NewRequest(http.MethodPost, "https://api.firecrawl.dev/v1/search", bytes.NewReader(body))
	if err != nil {
		return "search error"
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		if r.Log != nil {
			r.Log("[search]: request error:", err)
		}
		return "search failed"
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if r.Log != nil {
			r.Log("[search]: HTTP", resp.StatusCode, string(raw))
		}
		return fmt.Sprintf("search error: HTTP %d", resp.StatusCode)
	}

	var result searchResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return "search: failed to parse results"
	}

	if !result.Success || len(result.Data) == 0 {
		return "no results found"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("search results for: %s\n\n", query))
	for i, doc := range result.Data {
		title := doc.URL
		if doc.Title != "" {
			title = doc.Title
		}
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, title))
		if doc.URL != "" {
			sb.WriteString(fmt.Sprintf("    %s\n", doc.URL))
		}
		if doc.Description != "" {
			sb.WriteString(fmt.Sprintf("    %s\n", doc.Description))
		}
		sb.WriteString("\n")
	}

	if r.Log != nil {
		r.Log("[search]: found", len(result.Data), "results for:", query)
	}
	return sb.String()
}

type searchResponse struct {
	Success bool         `json:"success"`
	Data    []searchItem `json:"data"`
}

type searchItem struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
