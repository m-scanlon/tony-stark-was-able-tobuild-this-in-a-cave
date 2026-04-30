package reality

type Reality interface {
	ID() string
	Create(r *Relation) Reality
	Realize(r *Relation) string
}
