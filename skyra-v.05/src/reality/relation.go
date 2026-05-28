package reality

import (
	"fmt"
	"strings"
)

type Relation struct {
	id        string
	Origin    string
	ThreadID  string
	Impulse   string
	Providers map[string]Reality
	Fields    map[string]map[string]float64
	Depth     int
	Trace     []string
	Visited   map[string]bool
	Thoughts  []string
	Signal    float64
	MaxDepth  int
	Cancelled bool
}

func (r *Relation) ID() string {
	return r.id
}

func (r *Relation) Create(rel *Relation) Reality {
	return nil
}

func (r *Relation) Realize(rel *Relation) *Relation {
	return rel
}

func (r *Relation) Competence(rel *Relation) {
}

func (r *Relation) Observe(rel *Relation) {
	r.Depth++
	r.Trace = append(r.Trace, rel.ID())
}

func (r *Relation) Express(rel *Relation) *Relation {
	return rel
}

func (r *Relation) Activation(rel *Relation) float64 {
	return r.Signal
}

func (r *Relation) Reinforce() {
}

func (r *Relation) Decay() {
}

func Impress(origin, raw string) (*Relation, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("reality: empty input")
	}

	return &Relation{
		id:        NewID(),
		Origin:    origin,
		Impulse:   raw,
		Providers: make(map[string]Reality),
		Fields:    make(map[string]map[string]float64),
		Visited:   make(map[string]bool),
		Signal:    1.0,
		MaxDepth:  50,
	}, nil
}
