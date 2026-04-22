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

	parts := strings.SplitN(raw, "|", 2)
	if len(parts) != 2 {
		return Relation{}, fmt.Errorf("entity: missing | divider")
	}

	left := strings.TrimSpace(parts[0])
	tokens := strings.Fields(left)
	if len(tokens) < 2 {
		return Relation{}, fmt.Errorf("entity: expected at least protocol and target")
	}
	if tokens[0] != "skyra" {
		return Relation{}, fmt.Errorf("entity: must begin with skyra")
	}

	id := tokens[1]
	impulse := strings.Join(tokens[2:], " ")

	return Relation{
		ID:       id,
		Origin:   origin,
		ThreadID: threadID,
		Impulse:  impulse,
	}, nil
}
