package domain

import "errors"

var (
	ErrEmptyBeingName     = errors.New("domain: being name is required")
	ErrEmptyIdentity      = errors.New("domain: nature identity is required")
	ErrEmptyPurpose       = errors.New("domain: nature purpose is required")
	ErrNilBeing           = errors.New("domain: being is nil")
	ErrDuplicateBeingName = errors.New("domain: duplicate being name")
	ErrUnknownBeing       = errors.New("domain: unknown being")
	ErrUnknownPeer        = errors.New("domain: unknown peer")
	ErrInvalidImpulse     = errors.New("domain: invalid impulse")
)
