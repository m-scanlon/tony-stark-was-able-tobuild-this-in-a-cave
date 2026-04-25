package world

import (
	"skyra-v04/src/primitives/being"
	"skyra-v04/src/primitives/claude"
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
	"skyra-v04/src/primitives/opencode"
)

func (w *World) grow(r entity.Relation) {
	name, _ := meaning.Extract(r.Impulse, "~name", "grow")
	mediumName, _ := meaning.Extract(r.Impulse, "~medium", "grow")
	r.ID = name

	var e entity.Entity
	switch mediumName {
	case "claude":
		e = claude.Claude{}.Relate(r)
	case "opencode":
		e = opencode.OpenCode{}.Relate(r)
	default:
		e = being.Being{}.Relate(r)
	}
	w.EntityMap[e.ID()] = e
}
