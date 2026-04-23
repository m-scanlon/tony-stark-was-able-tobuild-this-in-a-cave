package thread

import (
	"fmt"
	"strconv"
	"strings"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/exchange"
	"skyra-v04/src/primitives/meaning"
)

type Thread struct {
	presentThread
	id            string
	About         string
	Because       string
	Active        bool
	Relationships map[RelationshipKey]exchange.Exchange
}

func (t Thread) ID() string   { return t.id }
func (t Thread) Name() string { return "thread" }

// Relate constructs a new Thread from a relation (zero-value semantics).
func (t Thread) Relate(r entity.Relation) entity.Entity {
	about, _ := meaning.Extract(r.Impulse, "~about", "thread")
	because, _ := meaning.Extract(r.Impulse, "~because", "thread")
	return Thread{
		id:            r.ThreadID,
		About:         about,
		Because:       because,
		Active:        true,
		Relationships: make(map[RelationshipKey]exchange.Exchange),
	}
}

// Append adds a relation to the exchange between two parties in this thread.
func (t Thread) Append(a, b string, r entity.Relation) Thread {
	key := NewRelationshipKey(a, b)
	newRel := make(map[RelationshipKey]exchange.Exchange, len(t.Relationships))
	for k, v := range t.Relationships {
		newRel[k] = v
	}
	newRel[key] = newRel[key].Append(r)
	t.Relationships = newRel
	return t
}

// ExchangeWith returns the exchange between beingID and peer in this thread.
func (t Thread) ExchangeWith(beingID, peer string) (exchange.Exchange, bool) {
	key := NewRelationshipKey(beingID, peer)
	ex, ok := t.Relationships[key]
	return ex, ok
}

// ExchangesFor returns formatted text of the exchange between beingID and peer,
// with each entry numbered [i].
func (t Thread) ExchangeBetween(a, b string) string {
	ex, ok := t.ExchangeWith(a, b)
	if !ok {
		return ""
	}
	var sb strings.Builder
	for i, rel := range ex.Relations {
		msg, err := meaning.Extract(rel.Impulse, "~say", "exchange", "|")
		if err != nil {
			msg = rel.Impulse
		}
		sb.WriteString(fmt.Sprintf("  [%d] %s: %s\n", i, rel.Origin, msg))
	}
	return sb.String()
}

// OtherExchangesFor returns a minimal summary of beingID's other exchanges in this
// thread — just peer names and entry counts. The being can ~ref to pull content.
func (t Thread) OtherExchangesFor(beingID, excludePeer string) string {
	var sb strings.Builder
	for key, ex := range t.Relationships {
		if key.A != beingID && key.B != beingID {
			continue
		}
		peer := key.A
		if key.A == beingID {
			peer = key.B
		}
		if peer == excludePeer {
			continue
		}
		sb.WriteString(fmt.Sprintf("  %s (%d entries)\n", peer, len(ex.Relations)))
	}
	return sb.String()
}

// ResolveRef parses a "<peer>:<start>-<end>" or "<peer>:<index>" ref, and
// returns the referenced slice of relations from beingID's exchange with <peer>.
func (t Thread) ResolveRef(beingID, ref string) []entity.Relation {
	colon := strings.Index(ref, ":")
	if colon == -1 {
		return nil
	}
	peer := strings.TrimSpace(ref[:colon])
	rng := strings.TrimSpace(ref[colon+1:])
	ex, ok := t.ExchangeWith(beingID, peer)
	if !ok {
		return nil
	}
	parts := strings.SplitN(rng, "-", 2)
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
	if start < 0 || start >= len(ex.Relations) {
		return nil
	}
	if end >= len(ex.Relations) {
		end = len(ex.Relations) - 1
	}
	return ex.Relations[start : end+1]
}

func (t Thread) DerivePresent(r entity.Relation) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("thread: %s\n", t.id))
	if t.About != "" {
		sb.WriteString("about: " + t.About + "\n")
	}
	if t.Because != "" {
		sb.WriteString("because: " + t.Because + "\n")
	}
	sb.WriteString(fmt.Sprintf("active: %v\n", t.Active))

	sb.WriteString("\nrelationships:\n")
	for key, ex := range t.Relationships {
		sb.WriteString(fmt.Sprintf("  %s <-> %s (%d entries)\n", key.A, key.B, len(ex.Relations)))
	}
	return sb.String()
}
