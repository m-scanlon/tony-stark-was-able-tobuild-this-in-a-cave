package reality

import "time"

type Reality interface {
	ID() string
	Core() *Base
	Create(r *Relation) Reality
	Realize(r *Relation) string
	Observe(r *Relation)
	Express(r *Relation) string
}

type Base struct {
	Weight        float64
	Usage         int
	LastUsed      time.Time
	Relationships map[string]Reality
	Expressors    map[string]Reality
}

func (b *Base) Core() *Base {
	return b
}

func (b *Base) Activation() float64 {
	return b.Weight
}
