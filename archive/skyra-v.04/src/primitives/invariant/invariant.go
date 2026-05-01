package invariant

import "skyra-v04/src/primitives/entity"

type Invariant struct{}

func (i Invariant) ID() string                          { return "" }
func (i Invariant) DerivePresent(_ entity.Relation) string { return "" }
func (i Invariant) Relate(r entity.Relation) entity.Entity { return i }
