package compress

import (
	"strings"
	"testing"
	"time"
)

func TestCompress_RespectsTokenBudgetAndSortsByScore(t *testing.T) {
	e := NewEngine(Options{
		MaxTokens:      20,
		MaxChunks:      3,
		MaxWordsPerHit: 20,
	})

	chunks := []Chunk{
		{ID: "low", Score: 0.1, Text: "low score text should likely drop"},
		{ID: "high", Score: 0.9, Text: "high score should be selected first"},
		{ID: "mid", Score: 0.5, Text: "mid score maybe selected based on budget"},
	}

	out := e.Compress(chunks)
	if len(out.Selected) == 0 {
		t.Fatalf("expected at least one selected chunk")
	}
	if out.Selected[0].ID != "high" {
		t.Fatalf("expected highest-score chunk first, got %q", out.Selected[0].ID)
	}
	if out.EstimatedTokens > 20 {
		t.Fatalf("token budget exceeded: %d", out.EstimatedTokens)
	}
}

func TestCompress_TrimsAndNormalizes(t *testing.T) {
	e := NewEngine(Options{
		MaxTokens:      200,
		MaxChunks:      1,
		MaxWordsPerHit: 4,
	})

	in := []Chunk{
		{
			ID:        "a",
			Score:     0.8,
			Timestamp: time.Now(),
			Text:      "  this   has\n many\tspaces and many words in a row  ",
		},
	}

	out := e.Compress(in)
	if len(out.Selected) != 1 {
		t.Fatalf("expected one selected chunk, got %d", len(out.Selected))
	}

	text := out.Selected[0].Text
	if !strings.HasSuffix(text, "...") {
		t.Fatalf("expected trimmed suffix, got %q", text)
	}
	if strings.Contains(text, "  ") || strings.Contains(text, "\n") || strings.Contains(text, "\t") {
		t.Fatalf("expected normalized whitespace, got %q", text)
	}
}

func TestCompress_EmptyInput(t *testing.T) {
	e := NewEngine(Options{})
	out := e.Compress(nil)
	if out.PromptBlock != "" {
		t.Fatalf("expected empty prompt block for empty input")
	}
}
