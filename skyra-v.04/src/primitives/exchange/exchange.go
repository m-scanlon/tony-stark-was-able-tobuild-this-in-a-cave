package exchange

import "skyra-v04/src/primitives/entity"

type Exchange struct {
	entity.PresentEntity
	Relations []entity.Relation
}

func (e Exchange) Append(r entity.Relation) Exchange {
	return Exchange{Relations: append(e.Relations, r)}
}

func (e Exchange) ID() string   { return "" }
func (e Exchange) Name() string { return "exchange" }

func (e Exchange) Relate(r entity.Relation) entity.Entity {
	return e.Append(r)
}
