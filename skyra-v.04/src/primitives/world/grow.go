package world

import (
	"skyra-v04/src/primitives/being"
	"skyra-v04/src/primitives/logos"
)

type Grow struct {
	logosMap map[string]logos.Logos
}

func (g *Grow) Relate(r logos.Relation) logos.Logos {
	b, _ := being.Being{}.Relate(r).(being.Being)
	g.logosMap[b.ID()] = b
	return g
}

func (g *Grow) ID() string   { return "grow" }
func (g *Grow) Name() string { return "grow" }
