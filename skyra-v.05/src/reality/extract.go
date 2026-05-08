package reality

import (
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
