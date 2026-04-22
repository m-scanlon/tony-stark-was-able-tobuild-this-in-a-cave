package world

import (
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/thread"
)

type World struct {
	LogosMap map[string]logos.Logos
	id       string
	name     string
}

func (w World) ID() string   { return w.id }
func (w World) Name() string { return w.name }

func (w World) Relate(r logos.Relation) logos.Logos {
	nodes := make(map[string]logos.Logos)
	newWorld := World{LogosMap: nodes}
	nodes["grow"] = &Grow{logosMap: nodes}
	nodes["start-thread"] = &thread.StartThread{LogosMap: nodes}
	nodes["continue-thread"] = &thread.ContinueThread{LogosMap: nodes}
	nodes["close-thread"] = &thread.CloseThread{LogosMap: nodes}
	nodes["parent"] = w
	return newWorld
}


