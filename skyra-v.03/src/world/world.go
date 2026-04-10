package world

import (
	"fmt"
	"strings"

	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/nature"
)

type World struct {
	beings map[string]*being.Being
}

func New() *World {
	return &World{
		beings: make(map[string]*being.Being),
	}
}

func (w *World) Grow(expression string) (*being.Being, error) {
	b, err := being.CreateBeing(expression)
	if err != nil {
		return nil, err
	}
	if err := w.Register(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (w *World) Register(b *being.Being) error {
	if b == nil {
		return ErrNilBeing
	}
	if err := b.Validate(); err != nil {
		return err
	}
	if _, exists := w.beings[b.Name]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateBeing, b.Name)
	}
	w.beings[b.Name] = b
	return nil
}

func (w *World) BeingByName(name string) (*being.Being, bool) {
	b, ok := w.beings[strings.TrimSpace(name)]
	return b, ok
}

func (w *World) Relate(leftName, rightName string) error {
	left, ok := w.beings[leftName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownBeing, leftName)
	}
	right, ok := w.beings[rightName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownBeing, rightName)
	}

	if err := seedPeer(left, right.Name, right.Nature); err != nil {
		return err
	}
	if err := seedPeer(right, left.Name, left.Nature); err != nil {
		return err
	}
	return nil
}

func SeedPeer(b *being.Being, peerName string, peerNature nature.Nature) error {
	return seedPeer(b, peerName, peerNature)
}

func seedPeer(b *being.Being, peerName string, peerNature nature.Nature) error {
	if _, exists := b.Peers[peerName]; exists {
		return nil
	}

	if b.Cognitive {
		peer, err := NewExchangeStack(peerName, peerNature)
		if err != nil {
			return err
		}
		return b.AttachPeer(peer)
	}

	peer, err := NewExternalDispatch(peerName, peerNature)
	if err != nil {
		return err
	}
	return b.AttachPeer(peer)
}
