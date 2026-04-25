package world

import (
	"fmt"
	"strings"

	being "skyra-v03/src/primitives/being"
)

type World struct {
	beings map[string]*being.Being
}

func New() *World {
	return &World{
		beings: make(map[string]*being.Being),
	}
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

func (w *World) FindOpenExchangeThread(beingName, peerName string) (string, bool) {
	b, ok := w.beings[strings.TrimSpace(beingName)]
	if !ok {
		return "", false
	}
	ch, ok := b.Peers[strings.TrimSpace(peerName)]
	if !ok {
		return "", false
	}
	em, ok := ch.(*ExchangeMap)
	if !ok {
		return "", false
	}
	for threadID := range em.exchanges {
		return threadID, true
	}
	return "", false
}

func (w *World) HasExchangeThread(beingName, peerName, threadID string) bool {
	if w == nil {
		return false
	}
	if strings.TrimSpace(threadID) == "" {
		return false
	}

	b, ok := w.beings[strings.TrimSpace(beingName)]
	if !ok {
		return false
	}
	peer, ok := b.Peers[strings.TrimSpace(peerName)]
	if !ok {
		return false
	}
	em, ok := peer.(*ExchangeMap)
	if !ok {
		return false
	}
	_, ok = em.ExchangeByThread(threadID)
	return ok
}
