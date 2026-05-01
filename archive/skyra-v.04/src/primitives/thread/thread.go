package thread

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/exchange"
	"skyra-v04/src/primitives/meaning"
)

func NewThreadID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

type Thread struct {
	id            string
	About         string
	Because       string
	Active        bool
	Relationships map[RelationshipKey]exchange.Exchange
}

func (t Thread) ID() string                                        { return t.id }
func (t Thread) Name() string                                      { return "thread" }
func (t Thread) ExchangeMap() map[RelationshipKey]exchange.Exchange { return t.Relationships }

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

// Append adds a relation to the exchange between two parties. If the exchange
// doesn't exist or is inactive, it's opened/reopened with r.Origin as the parent.
func (t Thread) Append(a, b string, r entity.Relation) Thread {
	key := NewRelationshipKey(a, b)
	newRel := make(map[RelationshipKey]exchange.Exchange, len(t.Relationships))
	for k, v := range t.Relationships {
		newRel[k] = v
	}
	ex := newRel[key]
	if !ex.Active {
		ex = ex.Open(r.Origin)
	}
	ex = ex.Append(r)
	newRel[key] = ex
	t.Relationships = newRel
	return t
}

// CloseExchange marks the exchange between a and b as inactive.
func (t Thread) CloseExchange(a, b string) Thread {
	key := NewRelationshipKey(a, b)
	newRel := make(map[RelationshipKey]exchange.Exchange, len(t.Relationships))
	for k, v := range t.Relationships {
		newRel[k] = v
	}
	if ex, ok := newRel[key]; ok {
		newRel[key] = ex.Close()
	}
	t.Relationships = newRel
	return t
}

// FindReturnTarget returns the parent of the most-recently-active exchange
// involving beingID where beingID is NOT the parent — i.e., an exchange the
// being was called INTO rather than one they opened. Returns "" if none.
func (t Thread) FindReturnTarget(beingID string) string {
	for key, ex := range t.Relationships {
		if !ex.Active {
			continue
		}
		if key.A != beingID && key.B != beingID {
			continue
		}
		if ex.Parent == beingID {
			continue
		}
		return ex.Parent
	}
	return ""
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
		sb.WriteString(fmt.Sprintf("  [%d] %s: %s\n", i, rel.Origin, rel.Impulse))
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

// ActiveExchangesFor returns all active exchanges involving beingID, annotated
// with role (current, waiting on you, you opened). currentPeer is the peer
// whose message the being is currently responding to.
func (t Thread) ActiveExchangesFor(beingID, currentPeer string) string {
	var sb strings.Builder
	for key, ex := range t.Relationships {
		if !ex.Active {
			continue
		}
		if key.A != beingID && key.B != beingID {
			continue
		}
		peer := key.A
		if key.A == beingID {
			peer = key.B
		}
		label := ""
		switch {
		case peer == currentPeer:
			label = "current"
		case ex.Parent != beingID:
			label = "waiting on your response"
		default:
			label = "you opened"
		}
		sb.WriteString(fmt.Sprintf("  %s — %s, %d entries\n", peer, label, len(ex.Relations)))
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
