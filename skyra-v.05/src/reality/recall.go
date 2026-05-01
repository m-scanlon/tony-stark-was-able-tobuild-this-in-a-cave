package reality

type Recall struct {
	id string
}

func (rc *Recall) ID() string { return rc.id }

func (rc *Recall) Create(r *Relation) Reality {
	return &Recall{id: "recall"}
}

func (rc *Recall) Realize(r *Relation) string {
	if r.Log != nil {
		r.Log("[recall]: no relevant memories")
	}
	return "no relevant memories"
}
