package world

import (
	"skyra-v04/src/primitives/being"
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

type Grow struct {
	presentWorld
	EntityMap map[string]entity.Entity
}

func (g *Grow) Relate(r entity.Relation) entity.Entity {
	name, _ := meaning.Extract(r.Impulse, "~name", "grow")
	r.ID = name
	b, _ := being.Being{}.Relate(r).(being.Being)
	g.EntityMap[b.ID()] = b
	return g
}

func (g *Grow) ID() string   { return "grow" }
func (g *Grow) Name() string { return "grow" }
