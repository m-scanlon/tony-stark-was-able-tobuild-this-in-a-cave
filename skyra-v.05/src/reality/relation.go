package reality

import (
	"fmt"
	"strings"
)

type Relation struct {
	Base
	RelationID string
	Origin     string
	ThreadID   string
	Impulse    string
	Log        func(args ...any)
	Exports    map[string]any
	Depth      int
	Trace      []string
	Visited    map[string]bool
	Thoughts   []string
	Signal     float64
	MaxDepth   int
}

func (r *Relation) ID() string {
	return r.RelationID
}

func (r *Relation) Create(rel *Relation) Reality {
	return nil
}

func (r *Relation) Observe(rel *Relation) {
	r.Depth++
	r.Trace = append(r.Trace, rel.ID())
}

func (r *Relation) ObserveReality(reality Reality) {
	r.Depth++
	r.Trace = append(r.Trace, reality.ID())
}

func (r *Relation) Export(key string, value any) {
	if r.Exports == nil {
		r.Exports = make(map[string]any)
	}
	r.Exports[key] = value
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
		Visited:   make(map[string]bool),
		Signal:    1.0,
		MaxDepth:  50,
	}, nil
}
