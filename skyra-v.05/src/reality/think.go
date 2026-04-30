package reality

type Think struct {
	id        string
	Realities map[string]Reality
}

func (t *Think) ID() string { return t.id }

func (t *Think) Create(r *Relation) Reality {
	return &Think{
		id:        "think",
		Realities: make(map[string]Reality),
	}
}

func (t *Think) Realize(r *Relation) string {
	r.Attach("think", t.Parse)
	return ""
}

func (t *Think) Parse() string {
	return ""
}
