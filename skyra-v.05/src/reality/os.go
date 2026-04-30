package reality

type OS struct {
	id        string
	Realities map[string]Reality
}

func (o *OS) ID() string { return o.id }

func (o *OS) Create(r *Relation) Reality {
	return &OS{
		id:        "os",
		Realities: make(map[string]Reality),
	}
}

func (o *OS) Realize(r *Relation) string {
	target, ok := o.Realities[r.ID]
	if !ok {
		return ""
	}
	return target.Realize(r)
}

func (o *OS) Parse() string {
	return ""
}
