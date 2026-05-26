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
	Depth      int
	Budget     float64
	Trace      []string
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
		Budget:    1.0,
	}, nil
}
