package world

import "errors"

var (
	ErrUnknownPeer    = errors.New("world: unknown peer")
	ErrInvalidImpulse = errors.New("world: invalid impulse")
	ErrNilBeing       = errors.New("world: being is nil")
	ErrUnknownBeing   = errors.New("world: unknown being")
	ErrDuplicateBeing = errors.New("world: duplicate being name")
)
