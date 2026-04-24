package world

import (
	"skyra-v04/src/primitives/being"
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

func (w *World) grow(r entity.Relation) {
	name, _ := meaning.Extract(r.Impulse, "~name", "grow")
	r.ID = name
	b, _ := being.Being{}.Relate(r).(being.Being)
	w.EntityMap[b.ID()] = b
}
