package reality

import (
	"fmt"
	"skyra-v05/src/debug"
	"strings"
)

type Self struct {
	id          string
	Realities   map[string]Reality
	Universe    *Universe
	Claimed     map[string]string
	Specialists map[string]*Self
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

	providers := make(map[string]Reality)
	if r.Realities != nil {
		for name, reality := range r.Realities {
			if _, ok := reality.(*Provider); ok {
				providers[name] = reality
			}
		}
	}

	desk := (&Desk{}).Create(&Relation{}).(*Desk)
	desk.Owner = r.ID
	self.Realities["desk"] = desk

	self.Claimed = make(map[string]string)
	self.Specialists = make(map[string]*Self)

	if len(providers) > 0 {
		mem := NewMemory(r.ID)
		if being, ok := self.Realities["being"].(Being); ok {
			mem.HomeDir = being.Home
			mem.Load()
			mem.SeedSkills("skills")
		}
		self.Realities["memory"] = mem

		ctx := &Context{
			id: "context", Owner: r.ID, Memory: mem, Providers: providers,
			Warm: make(map[string][]*MemNode),
			Claimed: self.Claimed, Specialists: self.Specialists,
			OnPromote: func(cluster *Cluster) { self.Promote(cluster, providers) },
		}
		self.Realities["context"] = ctx

		think := (&Think{}).Create(&Relation{}).(*Think)
		think.Providers = providers
		act := (&Act{}).Create(&Relation{}).(*Act)
		act.Providers = providers
		self.Realities["think"] = think
		self.Realities["act"] = act
	}

	return self
}

func (s *Self) Promote(cluster *Cluster, providers map[string]Reality) {
	log := func(args ...any) { debug.Being(s.id, "self", args...) }

	if s.Universe == nil {
		thread := &NewThread{
			id:       "thread-gate",
			Beings:   make(map[string]Reality),
			Access:   make(map[string]bool),
			Threads:  make(map[string]*Thread),
			Exchange: (&Exchange{}).Create(&Relation{}).(*Exchange),
			Devices:  make(map[string]Reality),
			ThinkOps: make(map[string]Reality),
			ActOps:   make(map[string]Reality),
		}
		s.Universe = &Universe{id: "inner-universe", Thread: thread}
		log("[self]: inner universe created")
	}

	var heaviest string
	var maxWeight float64
	for _, eid := range cluster.Entities {
		if e := s.Realities["memory"].(*Memory).Graph.GetEntity(eid); e != nil {
			if e.Weight > maxWeight {
				maxWeight = e.Weight
				heaviest = e.Name
			}
		}
	}

	name := heaviest + "-specialist"
	if _, exists := s.Specialists[name]; exists {
		name = heaviest + "-specialist-2"
	}

	var entityNames []string
	for _, eid := range cluster.Entities {
		if e := s.Realities["memory"].(*Memory).Graph.GetEntity(eid); e != nil {
			entityNames = append(entityNames, e.Name)
		}
	}

	impulse := fmt.Sprintf("~name %s\n~type llm\n~identity specialist processor for %s\n~purpose I hold deep understanding of %s for my parent\n~relationships %s",
		name, strings.Join(entityNames, ", "), strings.Join(entityNames, ", "), s.id)

	mem := s.Realities["memory"].(*Memory)

	specSelf := &Self{
		id:          name,
		Realities:   make(map[string]Reality),
		Claimed:     make(map[string]string),
		Specialists: make(map[string]*Self),
	}

	being := Being{}.Create(&Relation{Impulse: impulse}).(Being)
	specSelf.Realities["being"] = being

	specCtx := &Context{
		id: "context", Owner: name, Memory: mem, Providers: providers,
		Warm: make(map[string][]*MemNode),
		Scope: cluster.Entities, Claimed: specSelf.Claimed, Specialists: specSelf.Specialists,
		OnPromote: func(c *Cluster) { specSelf.Promote(c, providers) },
	}
	specSelf.Realities["context"] = specCtx
	specSelf.Realities["memory"] = mem

	think := (&Think{}).Create(&Relation{}).(*Think)
	think.Providers = providers
	think.Operators["retrieve-context"] = &RetrieveContext{}
	think.Operators["store-context"] = &StoreContext{}
	specSelf.Realities["think"] = think

	act := (&Act{}).Create(&Relation{}).(*Act)
	act.Providers = providers
	specSelf.Realities["act"] = act

	s.Universe.Thread.Beings[name] = specSelf
	s.Universe.Thread.Access[name] = true
	s.Universe.Thread.Access[s.id] = true
	s.Specialists[name] = specSelf

	for _, eid := range cluster.Entities {
		s.Claimed[eid] = name
		mem.Graph.Claimed[eid] = name
	}

	log("[self]: specialist promoted:", name, "entities:", entityNames, "density:", cluster.Density)
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

	thinkBacks := 0
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
				thinkBacks++
				if thinkBacks >= 3 {
					log("[self]: think-back budget exhausted")
					r.Origin = s.id
					return ""
				}
				log("[self]: think-back, re-entering")
				continue
			}

			r.Origin = s.id
			return result
		}

		return ""
	}
}
