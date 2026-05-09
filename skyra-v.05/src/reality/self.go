package reality

import (
	"skyra-v05/src/debug"
	"strings"
)

type Self struct {
	id        string
	Realities map[string]Reality
}

func (s *Self) ID() string { return s.id }

func (s *Self) Create(r *Relation) Reality {
	self := &Self{
		id:        r.ID,
		Realities: make(map[string]Reality),
	}

	if r.Impulse != "" {
		being := Being{}.Create(r).(Being)
		self.Realities["being"] = being
	}

	var llm Reality
	if r.Realities != nil {
		for _, reality := range r.Realities {
			if _, ok := reality.(*Provider); ok {
				llm = reality
				break
			}
		}
	}

	desk := (&Desk{}).Create(&Relation{}).(*Desk)
	desk.Owner = r.ID
	self.Realities["desk"] = desk

	if llm != nil {
		mem := NewMemory(r.ID)
		if being, ok := self.Realities["being"].(Being); ok {
			mem.HomeDir = being.Home
			mem.Load()
			mem.SeedSkills("skills")
		}
		self.Realities["memory"] = mem

		ctx := &Context{id: "context", Owner: r.ID, Memory: mem, LLM: llm, Warm: make(map[string][]*MemNode)}
		self.Realities["context"] = ctx

		think := (&Think{}).Create(&Relation{}).(*Think)
		think.LLM = llm
		act := (&Act{}).Create(&Relation{}).(*Act)
		act.LLM = llm
		self.Realities["think"] = think
		self.Realities["act"] = act
	}

	return self
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
			if being.Entrypoints != nil {
				snap.Entrypoints = being.Entrypoints
			}
			snap.Device = being.Device
			snap.Memories = snapshotMemories(being.Home)
			node.Children = append(node.Children, RealityNode{ID: s.id + "-being", Type: "Being", Children: []RealityNode{}})
		}

		if desk, ok := s.Realities["desk"].(*Desk); ok {
			desk.Realize(r)
			if ds, ok := r.Exports["desk:"+s.id]; ok {
				dsnap := ds.(DeskSnapshot)
				snap.Desk = &dsnap
				delete(r.Exports, "desk:"+s.id)
			}
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

	log := func(args ...any) { debug.Being(s.id, "self", args...) }
	log("[self]: realizing, origin:", r.Origin, "impulse:", truncate(r.Impulse, 60))

	if r.Realities == nil {
		r.Realities = make(map[string]Reality)
	}

	if being, ok := s.Realities["being"]; ok {
		r.Realities["being"] = being
	}

	if mem, ok := s.Realities["memory"]; ok {
		r.Realities["memory"] = mem
	}

	if ctx, ok := s.Realities["context"]; ok {
		r.Realities["context"] = ctx
		if c, ok := ctx.(*Context); ok {
			relationship := r.Origin
			c.Heat(relationship)
			parsed := c.Parse(relationship)
			if parsed != "" {
				log("[self]: memory context warm for", relationship)
				r.Attach("memory-context", func() string { return parsed })
			}
		}
	}

	if desk, ok := s.Realities["desk"]; ok {
		r.Realities["desk"] = desk
		desk.Realize(r)
	}

	for {
		if think, ok := s.Realities["think"]; ok {
			if t, ok := think.(*Think); ok {
				t.OuterOps = t.OuterOps[:0]
				if act, ok := s.Realities["act"]; ok {
					if a, ok := act.(*Act); ok {
						for name := range a.Operators {
							t.OuterOps = append(t.OuterOps, name)
						}
					}
				}
				if r.Realities != nil {
					for key := range r.Realities {
						if strings.HasPrefix(key, "act:") {
							t.OuterOps = append(t.OuterOps, strings.TrimPrefix(key, "act:"))
						}
					}
				}
			}
			log("[self]: firing think")
			inner := think.Realize(r)
			log("[self]: think →", truncate(inner, 80))

			r.Attach("inner", func() string { return inner })
		}

		if act, ok := s.Realities["act"]; ok {
			log("[self]: firing act")
			result := act.Realize(r)
			log("[self]: act →", truncate(result, 80))

			if r.ID == "_think" {
				log("[self]: think-back, re-entering")
				continue
			}

			r.Origin = s.id
			return result
		}

		return ""
	}
}
