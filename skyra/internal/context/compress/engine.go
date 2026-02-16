package compress

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

// Chunk is a retrieved context unit from memory/search layers.
type Chunk struct {
	ID        string
	ProjectID string
	Source    string
	Text      string
	Score     float64
	Timestamp time.Time
}

// Options controls compression behavior.
type Options struct {
	MaxTokens      int
	MaxChunks      int
	MaxWordsPerHit int
}

// Result is a compressed prompt block plus metadata for telemetry.
type Result struct {
	PromptBlock     string
	Selected        []Chunk
	Dropped         int
	EstimatedTokens int
}

var defaultOptions = Options{
	MaxTokens:      700,
	MaxChunks:      8,
	MaxWordsPerHit: 60,
}

// Engine provides deterministic context compression for prompt injection.
type Engine struct {
	opts Options
}

func NewEngine(opts Options) *Engine {
	cfg := defaultOptions
	if opts.MaxTokens > 0 {
		cfg.MaxTokens = opts.MaxTokens
	}
	if opts.MaxChunks > 0 {
		cfg.MaxChunks = opts.MaxChunks
	}
	if opts.MaxWordsPerHit > 0 {
		cfg.MaxWordsPerHit = opts.MaxWordsPerHit
	}
	return &Engine{opts: cfg}
}

// Compress ranks chunks, trims each chunk, and fits selected content into a token budget.
func (e *Engine) Compress(chunks []Chunk) Result {
	if len(chunks) == 0 {
		return Result{PromptBlock: ""}
	}

	ranked := make([]Chunk, 0, len(chunks))
	for _, c := range chunks {
		txt := normalize(c.Text)
		if txt == "" {
			continue
		}
		c.Text = txt
		ranked = append(ranked, c)
	}
	if len(ranked) == 0 {
		return Result{PromptBlock: ""}
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		// Higher score first, then newer first.
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].Timestamp.After(ranked[j].Timestamp)
		}
		return ranked[i].Score > ranked[j].Score
	})

	selected := make([]Chunk, 0, min(len(ranked), e.opts.MaxChunks))
	used := 0

	for _, c := range ranked {
		if len(selected) >= e.opts.MaxChunks {
			break
		}

		trimmed := trimWords(c.Text, e.opts.MaxWordsPerHit)
		toks := estimateTokens(trimmed)
		if toks <= 0 {
			continue
		}

		if used+toks > e.opts.MaxTokens {
			continue
		}

		c.Text = trimmed
		selected = append(selected, c)
		used += toks
	}

	lines := make([]string, 0, len(selected)+1)
	lines = append(lines, "Relevant context:")
	for i, c := range selected {
		src := c.Source
		if src == "" {
			src = "memory"
		}
		lines = append(lines, formatLine(i+1, src, c.Text))
	}

	return Result{
		PromptBlock:     strings.Join(lines, "\n"),
		Selected:        selected,
		Dropped:         len(ranked) - len(selected),
		EstimatedTokens: used,
	}
}

func formatLine(i int, source string, text string) string {
	return strings.Join([]string{
		"- [", itoa(i), "] ",
		source,
		": ",
		text,
	}, "")
}

func normalize(s string) string {
	parts := strings.Fields(strings.TrimSpace(s))
	return strings.Join(parts, " ")
}

func trimWords(s string, maxWords int) string {
	if maxWords <= 0 {
		return ""
	}
	words := strings.Fields(s)
	if len(words) <= maxWords {
		return s
	}
	return strings.Join(words[:maxWords], " ") + " ..."
}

func estimateTokens(s string) int {
	// Cheap approximation suitable for budgeting before final prompt assembly.
	words := len(strings.Fields(s))
	if words == 0 {
		return 0
	}
	return int(float64(words)*1.35) + 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func itoa(v int) string {
	return strconv.Itoa(v)
}
