package reality

type Plan struct {
	id string
}

func (p *Plan) ID() string { return p.id }

func (p *Plan) Create(r *Relation) Reality {
	return &Plan{id: "plan"}
}

func (p *Plan) Realize(r *Relation) string {
	if r.Log != nil {
		r.Log("[plan]: stub")
	}
	return ""
}
