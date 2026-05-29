package reality

import (
	"fmt"
	"strings"
	"sync"
)

type Relation struct {
	mu        sync.Mutex
	id        string
	Origin    string
	ThreadID  string
	Impulse   string
	Providers map[string]Reality
	Fields    map[string]map[string]float64 // memoryID → {bindingReality: activation}
	Depth     int
	Trace     []string
	Visited   map[string]bool
	LastSeen  Reality
	Thoughts  []string
	Signal    float64
	Load      int // accumulated content size in tokens/chars — checked on ascent
	MaxLoad   int
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
	r.mu.Lock()
	r.Depth++
	r.Trace = append(r.Trace, rel.ID())
	r.mu.Unlock()
}

func (r *Relation) Express(rel *Relation) *Relation {
	r.mu.Lock()
	spent := r.Signal < 0.1
	r.mu.Unlock()
	if spent {
		if host, ok := r.LastSeen.(*Skyra); ok {
			r.Deposit(host)
		}
	}
	return rel
}

func (r *Relation) Activation(rel *Relation) float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Signal
}

func (r *Relation) Deposit(host *Skyra) {
	r.mu.Lock()
	mem := host.Create(&Relation{}).(*Skyra)
	mem.Weight = r.Signal
	host.Relationships[mem.ID()] = mem
	r.Fields[mem.ID()] = map[string]float64{
		"thread:alignment": 1.0,
		"thread:strength":  1.0,
	}
	r.mu.Unlock()
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
