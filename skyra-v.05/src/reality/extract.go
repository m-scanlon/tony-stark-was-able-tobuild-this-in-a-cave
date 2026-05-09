package reality

import (
	"fmt"
	"strings"
	"unicode"
)

type Extractor struct {
	Known map[string]bool
}

func NewExtractor() *Extractor {
	return &Extractor{Known: make(map[string]bool)}
}

func (e *Extractor) Learn(name string) {
	e.Known[strings.ToLower(strings.TrimSpace(name))] = true
}

func (e *Extractor) Extract(text string) []string {
	var entities []string
	seen := map[string]bool{}

	lower := strings.ToLower(text)
	for name := range e.Known {
		if strings.Contains(lower, name) && !seen[name] {
			seen[name] = true
			entities = append(entities, name)
		}
	}

	words := strings.Fields(text)
	for i, word := range words {
		if i == 0 {
			continue
		}
		clean := strings.TrimRight(word, ".,;:!?()")
		if clean == "" {
			continue
		}
		runes := []rune(clean)
		if unicode.IsUpper(runes[0]) && len(runes) > 1 {
			name := strings.ToLower(clean)
			if !seen[name] && !isCommonWord(name) {
				seen[name] = true
				entities = append(entities, name)
			}
		}
	}

	return entities
}

type Resolver struct {
	Aliases map[string]string
}

func NewResolver() *Resolver {
	return &Resolver{Aliases: make(map[string]string)}
}

func (r *Resolver) AddAlias(alias, canonical string) {
	r.Aliases[normalizeEntity(alias)] = normalizeEntity(canonical)
}

func (r *Resolver) Resolve(name string) string {
	norm := normalizeEntity(name)
	if canonical, ok := r.Aliases[norm]; ok {
		return canonical
	}
	best := ""
	bestDist := 3
	for alias, canonical := range r.Aliases {
		d := levenshtein(norm, alias)
		if d > 0 && d < bestDist {
			bestDist = d
			best = canonical
		}
	}
	if best != "" {
		return best
	}
	return norm
}

func normalizeEntity(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func levenshtein(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(curr[j-1]+1, min(prev[j]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

var commonWords = map[string]bool{
	"i": true, "the": true, "a": true, "an": true, "is": true,
	"it": true, "we": true, "you": true, "they": true, "he": true,
	"she": true, "my": true, "your": true, "his": true, "her": true,
	"this": true, "that": true, "these": true, "those": true,
	"what": true, "when": true, "where": true, "who": true, "how": true,
	"but": true, "and": true, "or": true, "not": true, "no": true,
	"yes": true, "if": true, "then": true, "so": true, "too": true,
	"also": true, "just": true, "can": true, "will": true, "would": true,
	"should": true, "could": true, "do": true, "did": true, "does": true,
	"have": true, "has": true, "had": true, "been": true, "being": true,
	"are": true, "was": true, "were": true, "am": true, "be": true,
	"for": true, "with": true, "from": true, "into": true, "about": true,
	"after": true, "before": true, "between": true, "under": true, "over": true,
	"here": true, "there": true, "now": true, "still": true,
}

func isCommonWord(w string) bool {
	return commonWords[w]
}

func Extract(expression, token, name string, delimiters ...string) (string, error) {
	idx := strings.Index(expression, token)
	if idx == -1 {
		return "", fmt.Errorf("%s: token %q not found in expression", name, token)
	}

	rest := strings.TrimSpace(expression[idx+len(token):])
	if rest == "" {
		return "", fmt.Errorf("%s: no value after token %q", name, token)
	}

	if len(delimiters) == 0 {
		delimiters = []string{"~", "|"}
	}

	end := len(rest)
	for _, delim := range delimiters {
		if i := strings.Index(rest, delim); i != -1 && i < end {
			end = i
		}
	}

	value := strings.TrimSpace(rest[:end])
	if value == "" {
		return "", fmt.Errorf("%s: empty value for token %q", name, token)
	}
	return value, nil
}

func ExtractTag(text, tag string) (string, error) {
	open := "<" + tag + ">"
	close := "</" + tag + ">"
	start := strings.Index(text, open)
	if start == -1 {
		return "", fmt.Errorf("tag %q not found", tag)
	}
	after := text[start+len(open):]
	end := strings.Index(after, close)
	if end == -1 {
		return "", fmt.Errorf("tag %q not closed", tag)
	}
	return strings.TrimSpace(after[:end]), nil
}

func StripTag(text, tag string) string {
	open := "<" + tag + ">"
	close := "</" + tag + ">"
	start := strings.Index(text, open)
	if start == -1 {
		return text
	}
	after := text[start+len(open):]
	end := strings.Index(after, close)
	if end == -1 {
		return text
	}
	before := strings.TrimSpace(text[:start])
	rest := strings.TrimSpace(after[end+len(close):])
	if before == "" {
		return rest
	}
	if rest == "" {
		return before
	}
	return before + " " + rest
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
