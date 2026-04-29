package world

import (
	"skyra-v05/src/reality"
	"skyra-v05/src/reality/being"
)

type BeingWorld struct {
	World
	Pathos being.Being
	device Device
}

func NewBeingWorld(pathos being.Being, device Device) *BeingWorld {
	return &BeingWorld{
		World: World{
			id:        pathos.ID(),
			name:      pathos.Name(),
			Realities: make(map[string]reality.Reality),
		},
		Pathos: pathos,
		device: device,
	}
}

func (bw *BeingWorld) Realize(r reality.Relation) string {
	present := bw.Pathos.Realize(r)

	return bw.device.Realize(reality.Relation{
		ID:       r.ID,
		Origin:   r.Origin,
		ThreadID: r.ThreadID,
		Impulse:  present + "\n" + r.Impulse,
	})
}
