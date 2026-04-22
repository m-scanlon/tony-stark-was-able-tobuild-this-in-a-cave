package impression

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

type Impression struct {
	entity.PresentEntity
	Value string
}

func (i Impression) ID() string   { return i.Value }
func (i Impression) Name() string { return "impression" }

func (i Impression) Relate(r entity.Relation) entity.Entity {
	value, _ := meaning.Extract(r.Impulse, "~impression", "impression")
	return Impression{Value: value}
}
