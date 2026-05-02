package reality

import "skyra-v05/src/debug"

type Self struct {
	id        string
	Realities map[string]Reality
}

func (s *Self) ID() string { return s.id }

func (s *Self) Create(r *Relation) Reality {
	return &Self{
		id:        r.ID,
		Realities: make(map[string]Reality),
	}
}

func (s *Self) Realize(r *Relation) string {
	if r.Collecting {
		node := RealityNode{ID: s.id, Type: "Self", Children: []RealityNode{}}
		snap := BeingSnapshot{
			Name: s.id, Type: "llm", Status: "idle",
			Peers:    []string{},
			Memories: MemorySnapshot{Items: []MemoryItem{}, Skills: []SkillItem{}},
		}

		if being, ok := s.Realities["being"].(Being); ok {
			snap.Identity = being.Identity
			snap.Purpose = being.Purpose
			if being.Relationships != nil {
				snap.Peers = being.Relationships
			}
			snap.Device = being.Device
			snap.Memories = snapshotMemories(being.Home)
			node.Children = append(node.Children, RealityNode{ID: s.id + "-being", Type: "Being", Children: []RealityNode{}})
		}

		layers := &LayersSnapshot{}
		if think, ok := s.Realities["think"].(*Think); ok {
			think.Realize(r)
			if ts, ok := r.Exports["think"]; ok {
				layers.Think = ts.(ThinkSnapshot)
				delete(r.Exports, "think")
			}
			if tn, ok := r.Exports["node:think"]; ok {
				node.Children = append(node.Children, tn.(RealityNode))
				delete(r.Exports, "node:think")
			}
		}
		if act, ok := s.Realities["act"].(*Act); ok {
			act.Realize(r)
			if as, ok := r.Exports["act"]; ok {
				layers.Act = as.(ActSnapshot)
				delete(r.Exports, "act")
			}
			if an, ok := r.Exports["node:act"]; ok {
				node.Children = append(node.Children, an.(RealityNode))
				delete(r.Exports, "node:act")
			}
		}
		snap.Layers = layers

		r.Export("being:"+s.id, snap)
		r.Export("node:"+s.id, node)
		return ""
	}

	debug.Log("[self]: realizing", s.id)

	if being, ok := s.Realities["being"]; ok {
		debug.Log("[self]: passing being to relation")
		if r.Realities == nil {
			r.Realities = make(map[string]Reality)
		}
		r.Realities["being"] = being
	}

	outerParsers := r.Parsers
	r.Parsers = make(map[string]Parser)

	if think, ok := s.Realities["think"]; ok {
		if t, ok := think.(*Think); ok {
			if act, ok := s.Realities["act"]; ok {
				if a, ok := act.(*Act); ok {
					t.OuterOps = t.OuterOps[:0]
					for name := range a.Operators {
						t.OuterOps = append(t.OuterOps, name)
					}
				}
			}
		}
		debug.Log("[self]: firing think")
		inner := think.Realize(r)
		debug.Log("[self]: think returned:", inner)

		r.Parsers = make(map[string]Parser)
		r.Attach("inner", func() string { return inner })
	}

	if act, ok := s.Realities["act"]; ok {
		debug.Log("[self]: firing act")
		for name, parser := range outerParsers {
			r.Attach(name, parser)
		}
		result := act.Realize(r)
		debug.Log("[self]: act returned:", result)
		r.Origin = s.id
		return result
	}

	return ""
}
