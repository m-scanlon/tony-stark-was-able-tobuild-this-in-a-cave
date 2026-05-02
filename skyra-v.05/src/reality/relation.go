package reality

import (
	"fmt"
	"strings"
)

type Parser func() string

type Relation struct {
	ID         string
	Origin     string
	ThreadID   string
	Impulse    string
	Parsers    map[string]Parser
	Realities  map[string]Reality
	Log        func(args ...any)
	Collecting bool
	Exports    map[string]any
}

func (r *Relation) Export(key string, value any) {
	if r.Exports == nil {
		r.Exports = make(map[string]any)
	}
	r.Exports[key] = value
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

type Error struct {
	id      string
	Message string
}

func (e *Error) ID() string                  { return e.id }
func (e *Error) Create(r *Relation) Reality   { return e }
func (e *Error) Realize(r *Relation) string   { return e.Message }

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
	rest := response
	for {
		openStart := strings.Index(rest, "<")
		if openStart == -1 {
			break
		}
		openEnd := strings.Index(rest[openStart:], ">")
		if openEnd == -1 {
			break
		}
		target := rest[openStart+1 : openStart+openEnd]
		if target == "" || strings.ContainsAny(target, " /!?") {
			rest = rest[openStart+openEnd+1:]
			continue
		}
		closeTag := "</" + target + ">"
		after := rest[openStart+openEnd+1:]
		closeIdx := strings.Index(after, closeTag)
		var message string
		if closeIdx != -1 {
			message = strings.TrimSpace(after[:closeIdx])
			rest = after[closeIdx+len(closeTag):]
		} else {
			message = strings.TrimSpace(after)
			rest = ""
		}
		if message != "" {
			relations = append(relations, &Relation{
				Origin:    origin,
				ID:        strings.ToLower(target),
				Impulse:   message,
				Parsers:   make(map[string]Parser),
				Realities: make(map[string]Reality),
			})
		}
		if rest == "" {
			break
		}
	}
	return relations
}
