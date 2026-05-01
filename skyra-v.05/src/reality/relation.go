package reality

import (
	"fmt"
	"strings"
)

type Parser func() string

type Relation struct {
	ID        string
	Origin    string
	ThreadID  string
	Impulse   string
	Parsers   map[string]Parser
	Realities map[string]Reality
	Log       func(args ...any)
}

func (r *Relation) Attach(name string, parser Parser) {
	if r.Parsers == nil {
		r.Parsers = make(map[string]Parser)
	}
	r.Parsers[name] = parser
}

func Impress(origin, raw string) (*Relation, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("reality: empty input")
	}

	return &Relation{
		Origin:    origin,
		Impulse:   raw,
		Parsers:   make(map[string]Parser),
		Realities: make(map[string]Reality),
	}, nil
}

func (r *Relation) Peel() (string, string) {
	tokens := strings.Fields(r.Impulse)
	if len(tokens) == 0 {
		return "", ""
	}
	token := strings.ToLower(strings.TrimRight(tokens[0], ",:;."))
	rest := strings.Join(tokens[1:], " ")
	return token, rest
}

func ParseResponse(origin, response string) []*Relation {
	var relations []*Relation
	parts := strings.Split(response, "<>")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		target, message, found := strings.Cut(part, "|")
		if !found {
			continue
		}
		target = strings.TrimSpace(target)
		message = strings.TrimSpace(message)
		if target == "" || message == "" {
			continue
		}
		relations = append(relations, &Relation{
			Origin:    origin,
			ID:        strings.ToLower(target),
			Impulse:   message,
			Parsers:   make(map[string]Parser),
			Realities: make(map[string]Reality),
		})
	}
	return relations
}
