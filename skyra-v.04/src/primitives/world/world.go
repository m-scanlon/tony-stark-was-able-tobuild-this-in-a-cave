package world

import (
	"skyra-v04/src/primitives/adapter"
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
	l := make(map[string]logos.Logos)
	newWorld := World{LogosMap: l}
	l["grow"] = &Grow{LogosMap: l}
	l["spawn"] = &adapter.SpawnLogos{LogosMap: l}
	l["learn"] = &adapter.LearnLogos{LogosMap: l}
	l["start-thread"] = &thread.StartThread{LogosMap: l}
	l["continue-thread"] = &thread.ContinueThread{LogosMap: l}
	l["close-thread"] = &thread.CloseThread{LogosMap: l}
	l["parent"] = w
	l["read"] = adapter.ReadLogos{}
	l["find"] = adapter.FindLogos{}
	l["grep"] = adapter.GrepLogos{}
	l["write"] = adapter.WriteLogos{}
	l["append"] = adapter.AppendLogos{}
	l["delete"] = adapter.DeleteLogos{}
	l["move"] = adapter.MoveLogos{}
	l["mkdir"] = adapter.MkdirLogos{}
	l["list"] = adapter.ListLogos{}
	return newWorld
}


