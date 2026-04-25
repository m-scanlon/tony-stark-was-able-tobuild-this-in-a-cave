package world

import (
	"fmt"
	"strings"

	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/meaning"
	"skyra-v03/src/primitives/nature"
)

func (w *World) Grow(expression string) (*being.Being, error) {
	name, err := meaning.Extract(expression, "~name", "grow")
	if err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)

	if existing, ok := w.beings[name]; ok {
		return existing, w.seedRelationships(existing, expression)
	}

	b, err := being.CreateBeing(expression)
	if err != nil {
		return nil, err
	}
	if err := w.Register(b); err != nil {
		return nil, err
	}
	return b, w.seedRelationships(b, expression)
}

func (w *World) seedRelationships(b *being.Being, expression string) error {
	value, err := meaning.Extract(expression, "~relationships", "grow")
	if err != nil {
		return nil
	}
	for _, peerName := range strings.Split(value, ",") {
		peerName = strings.TrimSpace(peerName)
		if peerName == "" {
			continue
		}
		peer, ok := w.beings[peerName]
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnknownBeing, peerName)
		}
		if err := seedPeer(b, peer.Name, peer.Nature); err != nil {
			return err
		}
		if err := seedPeer(peer, b.Name, b.Nature); err != nil {
			return err
		}
		if lang, err := meaning.ExtractToEnd(expression, "~language-"+peerName, "grow"); err == nil {
			if ch, ok := b.Peers[peerName]; ok {
				if setter, ok := ch.(interface{ SetCallableLanguage(string) }); ok {
					setter.SetCallableLanguage(lang)
				}
			}
		}
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
		peer, err := NewExchangeMap(peerName, peerNature)
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
