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
		debug.Log("[self]: firing being")
		being.Realize(r)
	}

	if think, ok := s.Realities["think"]; ok {
		debug.Log("[self]: firing think")
		inner := think.Realize(r)
		debug.Log("[self]: think returned:", inner)
		r.Attach("inner", func() string { return inner })
	}

	if act, ok := s.Realities["act"]; ok {
		debug.Log("[self]: firing act")
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
