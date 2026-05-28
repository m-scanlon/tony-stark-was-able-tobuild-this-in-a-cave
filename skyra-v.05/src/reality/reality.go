package reality

import "math"

type Reality interface {
	ID() string
	Core() *Base
	Create(r *Relation) Reality
	Realize(r *Relation) string
	Observe(r *Relation)
	Express(r *Relation) string
}

type Base struct {
	Weight         float64
	TraversalCount int
	LastTraversed  int
	Alpha          float64
	Threads        map[string]bool
	Relationships  map[string]Reality
	Expressors     map[string]Reality
	Providers      map[string]Reality
}

func (b *Base) ID() string {
	return ""
}

func (b *Base) Core() *Base {
	return b
}

func (b *Base) Create(r *Relation) Reality {
	return nil
}

func (b *Base) Realize(r *Relation) string {
	if r.Visited[b.ID()] {
		return ""
	}
	r.Visited[b.ID()] = true

	// mutual observation — both are transformed
	r.Observe(b)
	b.Observe(r)

	// providers fire — think pass
	for _, p := range b.Providers {
		if thought := p.Realize(r); thought != "" {
			r.Thoughts = append(r.Thoughts, thought)
		}
	}

	// descent — deeper into relationships
	for _, rel := range b.Relationships {
		rel.Realize(r)
	}

	// ascent — compression and action formation
	result := b.Express(r)

	for _, e := range b.Expressors {
		e.Realize(r)
	}

	// return path — weight updates
	b.TraversalCount++
	b.LastTraversed = b.TraversalCount
	b.Reinforce()

	return result
}

func (b *Base) Observe(r *Relation) {}

func (b *Base) Express(r *Relation) string {
	return ""
}

func (b *Base) Recency(beingTraversalCount int) float64 {
	delta := beingTraversalCount - b.LastTraversed
	if delta <= 0 {
		return 1.0
	}
	return math.Pow(float64(delta), -0.5)
}

func (b *Base) ThreadAlignment(rel *Relation) float64 {
	if b.Threads == nil {
		return 1.0
	}
	if b.Threads[rel.ThreadID] {
		return 1.0
	}
	return 0.0
}

func (b *Base) Activation(beingTraversalCount int, rel *Relation) float64 {
	return b.Weight * b.Recency(beingTraversalCount) * b.ThreadAlignment(rel)
}

func (b *Base) Reinforce() {
	b.Weight = b.Alpha*1.0 + (1-b.Alpha)*b.Weight
}

func (b *Base) Decay() {
	b.Weight = (1 - b.Alpha) * b.Weight
}
