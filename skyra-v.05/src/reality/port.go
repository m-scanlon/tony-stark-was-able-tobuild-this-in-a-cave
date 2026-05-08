package reality

type Port interface {
	Reality
	Render(r *Relation) (string, string)
}
