package entity

import (
	"fmt"
	"strings"
)

type Relation struct {
	ID       string
	Origin   string
	ThreadID string
	Impulse  string
}

// Entity is the universal interface. Every node in the runtime — being, world, operator, adapter — is an Entity.
// Nodes relate to each other by passing a Relation. Nothing else is required.
type Entity interface {
	Relate(r Relation) Entity
	ID() string
	DerivePresent(r Relation) string
}

func Impress(origin, threadID, raw string) (Relation, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Relation{}, fmt.Errorf("entity: empty input")
	}

	tokens := strings.Fields(raw)
	if len(tokens) < 2 {
		return Relation{}, fmt.Errorf("entity: expected at least target and message")
	}

	id := strings.ToLower(strings.TrimRight(tokens[0], ",:;."))
	impulse := strings.Join(tokens[1:], " ")

	return Relation{
		ID:       id,
		Origin:   origin,
		ThreadID: threadID,
		Impulse:  impulse,
	}, nil
}
