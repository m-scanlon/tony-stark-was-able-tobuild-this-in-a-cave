package world

import (
	"skyra-v04/src/primitives/being"
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

type Grow struct {
	LogosMap map[string]logos.Logos
}

func (g *Grow) Relate(r logos.Relation) logos.Logos {
	name, _ := meaning.Extract(r.Impulse, "~name", "grow")
	r.ID = name
	b, _ := being.Being{}.Relate(r).(being.Being)
	g.LogosMap[b.ID()] = b
	return g
}

func (g *Grow) ID() string   { return "grow" }
func (g *Grow) Name() string { return "grow" }
