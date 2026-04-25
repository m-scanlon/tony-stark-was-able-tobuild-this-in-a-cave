package world

import (
	"fmt"
	"strconv"
	"strings"

	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/meaning"
)

type StartExchangeResult struct {
	PeerName      string
	ThreadID      string
	About         string
	Because       string
	ExpressionRef string
	Said          string
}

type CloseExchangeResult struct {
	PeerName string
}


func (w *World) StartExchange(expression string) (StartExchangeResult, error) {
	peerName, err := meaning.Extract(expression, "~with", "start-exchange")
	if err != nil {
		return StartExchangeResult{}, fmt.Errorf("start-exchange: %w", err)
	}
	about, err := meaning.Extract(expression, "~about", "start-exchange")
	if err != nil {
		return StartExchangeResult{}, fmt.Errorf("start-exchange: %w", err)
	}
	because, err := meaning.Extract(expression, "~because", "start-exchange")
	if err != nil {
		return StartExchangeResult{}, fmt.Errorf("start-exchange: %w", err)
	}
	said, err := meaning.ExtractToEnd(expression, "~say", "start-exchange")
	if err != nil {
		return StartExchangeResult{}, fmt.Errorf("start-exchange: %w", err)
	}
	ref, _ := meaning.Extract(expression, "~expression-reference", "start-exchange")
	return StartExchangeResult{PeerName: peerName, ThreadID: about, About: about, Because: because, ExpressionRef: ref, Said: said}, nil
}

func (w *World) OpenExchange(openerName, peerName, threadID, about, because string, contextEntries []being.ExchangeEntry) error {
	opener, ok := w.beings[openerName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownBeing, openerName)
	}
	ch, ok := opener.Peers[peerName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownPeer, peerName)
	}
	em, ok := ch.(*ExchangeMap)
	if !ok {
		return fmt.Errorf("world: peer %s does not support exchange lifecycle", peerName)
	}
	return em.OpenExchange(threadID, about, because, contextEntries)
}

func (w *World) ResolveExpressionRef(beingName, threadID, ref string) []being.ExchangeEntry {
	if ref == "" {
		return nil
	}
	b, ok := w.beings[beingName]
	if !ok {
		return nil
	}
	parts := strings.SplitN(ref, "-", 2)
	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil
	}
	end := start
	if len(parts) == 2 {
		end, err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil
		}
	}
	for _, ch := range b.Peers {
		em, ok := ch.(*ExchangeMap)
		if !ok {
			continue
		}
		thread, ok := em.exchanges[threadID]
		if !ok {
			continue
		}
		if start >= len(thread.Entries) {
			return nil
		}
		if end >= len(thread.Entries) {
			end = len(thread.Entries) - 1
		}
		return thread.Entries[start : end+1]
	}
	return nil
}

func (w *World) ParseCloseExchange(expression string) (CloseExchangeResult, error) {
	peerName, err := meaning.Extract(expression, "~with", "close-exchange")
	if err != nil {
		return CloseExchangeResult{}, fmt.Errorf("close-exchange: %w", err)
	}
	return CloseExchangeResult{PeerName: peerName}, nil
}

func (w *World) CloseExchange(openerName, peerName, threadID string) error {
	opener, ok := w.beings[openerName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownBeing, openerName)
	}
	ch, ok := opener.Peers[peerName]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownPeer, peerName)
	}
	em, ok := ch.(*ExchangeMap)
	if !ok {
		return fmt.Errorf("world: peer %s does not support exchange lifecycle", peerName)
	}
	return em.CloseExchange(threadID)
}
