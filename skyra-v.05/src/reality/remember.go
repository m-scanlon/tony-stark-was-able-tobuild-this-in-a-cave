package reality

type Remember struct {
	id string
}

func (rm *Remember) ID() string { return rm.id }

func (rm *Remember) Create(r *Relation) Reality {
	return &Remember{id: "remember"}
}

func (rm *Remember) Realize(r *Relation) string {
	if r.Log != nil {
		r.Log("[remember]: action not available")
	}
	return "action not available at this time"
}
