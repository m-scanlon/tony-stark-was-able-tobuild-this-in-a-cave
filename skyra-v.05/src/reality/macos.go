package reality

type MacOS struct {
	id         string
	Components map[string]Reality
}

func (m *MacOS) ID() string { return m.id }

func (m *MacOS) Create(r *Relation) Reality {
	id := r.ID
	if id == "" {
		id = "macos"
	}
	return &MacOS{
		id:         id,
		Components: make(map[string]Reality),
	}
}

func (m *MacOS) Component(name string) Reality {
	return m.Components[name]
}

func (m *MacOS) Realize(r *Relation) string {
	if r.Collecting {
		node := RealityNode{ID: m.id, Type: "MacOS", Children: []RealityNode{}}
		for name, comp := range m.Components {
			node.Children = append(node.Children, RealityNode{
				ID: name, Type: capitalizeType(name), Children: []RealityNode{},
			})
			comp.Realize(r)
		}
		r.Export("node:"+m.id, node)
		return ""
	}

	target, ok := r.Parsers["device-target"]
	if ok {
		name := target()
		if comp, exists := m.Components[name]; exists {
			return comp.Realize(r)
		}
	}

	if term, exists := m.Components["terminal"]; exists {
		return term.Realize(r)
	}

	return ""
}
