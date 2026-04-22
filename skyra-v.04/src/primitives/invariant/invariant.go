package invariant

import "skyra-v04/src/primitives/entity"

type Invariant struct {
	entity.PresentEntity
}

func (i Invariant) ID() string { return "" }

func (i Invariant) Relate(r entity.Relation) entity.Entity { return i }
