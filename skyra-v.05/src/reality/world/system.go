package world

import (
	"skyra-v05/src/reality"
)

type System struct {
	World
}

func NewSystem(physics *Physics) *System {
	return &System{
		World: World{
			id:        "system",
			name:      "system",
			Realities: make(map[string]reality.Reality),
			physics:   physics,
		},
	}
}

func (s *System) Realize(r reality.Relation) string {
	context := s.physics.Realize(r)

	target, ok := s.Realities[r.ID]
	if !ok {
		return ""
	}

	r.Impulse = context + "\n" + r.Impulse
	return target.Realize(r)
}
