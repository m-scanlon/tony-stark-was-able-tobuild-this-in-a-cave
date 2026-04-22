package language

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

type Language struct {
	entity.PresentEntity
	Value string
}

func (l Language) ID() string   { return l.Value }
func (l Language) Name() string { return "language" }

func (l Language) Relate(r entity.Relation) entity.Entity {
	value, _ := meaning.Extract(r.Impulse, "~language", "language")
	return Language{Value: value}
}
