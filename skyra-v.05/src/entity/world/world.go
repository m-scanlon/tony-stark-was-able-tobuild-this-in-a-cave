// World is the base container. It holds a hashmap of entities. World types
// (system, being, llm) embed this and provide their own DerivePresent.
package world

import (
	"skyra-v05/src/entity"
)

type World struct {
	id        string
	name      string
	Entities  map[string]entity.Entity
}

func New() *World {
	return &World{
		Entities: make(map[string]entity.Entity),
	}
}

func (w *World) ID() string   { return w.id }
func (w *World) Name() string { return w.name }

func (w *World) Create(r entity.Relation) entity.Entity {
	return w
}

func (w *World) DerivePresent(r entity.Relation) string {
	return ""
}
