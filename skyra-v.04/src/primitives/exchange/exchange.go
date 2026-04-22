package exchange

import "skyra-v04/src/primitives/entity"

type Exchange struct {
	entity.PresentEntity
	entries []string
}

func (e Exchange) Relate(r entity.Relation) entity.Entity {
	return Exchange{entries: append(e.entries, r.Impulse)}
}

func (e Exchange) ID() string      { return "" }
func (e Exchange) Name() string    { return "exchange" }
func (e Exchange) Entries() []string { return e.entries }
