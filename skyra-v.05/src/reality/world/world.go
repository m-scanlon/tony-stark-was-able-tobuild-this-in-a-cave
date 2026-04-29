package world

import (
	"skyra-v05/src/reality"
)

type World struct {
	id        string
	name      string
	Realities map[string]reality.Reality
	physics   *Physics
}

func (w World) ID() string { return w.id }

func (w World) Create(r reality.Relation) reality.Reality {
	return w
}

func (w World) Realize(r reality.Relation) string {
	target, ok := w.Realities[r.ID]
	if !ok {
		return ""
	}
	return target.Realize(r)
}
