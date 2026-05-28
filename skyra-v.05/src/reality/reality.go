package reality

import "math"

type Reality interface {
	ID() string
	Create(r *Relation) Reality
	Realize(r *Relation) *Relation
	Competence(r *Relation)
	Observe(r *Relation)
	Express(r *Relation) *Relation
	Activation(rel *Relation) float64
	Reinforce()
	Decay()
}

func Recency(traversalCount, lastTraversed int) float64 {
	delta := traversalCount - lastTraversed
	if delta <= 0 {
		return 1.0
	}
	return math.Pow(float64(delta), -0.5)
}
