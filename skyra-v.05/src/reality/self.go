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

func (s *Self) Parse() string {
	return ""
}
