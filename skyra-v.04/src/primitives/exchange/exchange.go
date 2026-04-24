package exchange

import "skyra-v04/src/primitives/entity"

// Exchange is the history between two beings in a thread, plus metadata about
// who currently "owns" the exchange (Parent) and whether it's still open (Active).
type Exchange struct {
	Parent    string // the being who most recently opened or reopened this exchange
	Active    bool   // whether this exchange is currently open
	Relations []entity.Relation
}

func (e Exchange) Append(r entity.Relation) Exchange {
	return Exchange{
		Parent:    e.Parent,
		Active:    e.Active,
		Relations: append(e.Relations, r),
	}
}

// Open marks the exchange active and sets the given being as its parent.
// Called on first creation and every time an inactive exchange reopens.
func (e Exchange) Open(parent string) Exchange {
	return Exchange{
		Parent:    parent,
		Active:    true,
		Relations: e.Relations,
	}
}

// Close marks the exchange inactive without touching its parent or history.
func (e Exchange) Close() Exchange {
	return Exchange{
		Parent:    e.Parent,
		Active:    false,
		Relations: e.Relations,
	}
}

func (e Exchange) ID() string                          { return "" }
func (e Exchange) Name() string                        { return "exchange" }
func (e Exchange) DerivePresent(_ entity.Relation) string { return "" }

func (e Exchange) Relate(r entity.Relation) entity.Entity {
	return e.Append(r)
}
