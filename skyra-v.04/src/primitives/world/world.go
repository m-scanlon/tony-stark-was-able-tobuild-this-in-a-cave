package world

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/thread"
)

type World struct {
	presentWorld
	EntityMap map[string]entity.Entity
	id       string
	name     string
}

func (w World) ID() string   { return w.id }
func (w World) Name() string { return w.name }

func (w World) Relate(r entity.Relation) entity.Entity {
	l := make(map[string]entity.Entity)
	newWorld := World{EntityMap: l}
	l["grow"] = &Grow{EntityMap: l}
	l["start-thread"] = &thread.StartThread{EntityMap: l}
	l["continue-thread"] = &thread.ContinueThread{EntityMap: l}
	l["close-thread"] = &thread.CloseThread{EntityMap: l}
	l["parent"] = w
	return newWorld
}


