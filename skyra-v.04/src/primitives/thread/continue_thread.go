package thread

import (
	"fmt"
	"strings"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
	"skyra-v04/src/primitives/medium"
	"skyra-v04/src/primitives/operator"
)

var _ operator.IOperator = (*ContinueThread)(nil)

type inferrable interface {
	entity.Entity
	Name() string
	Medium() medium.Medium
}

type ContinueThread struct {
	presentThread
	EntityMap map[string]entity.Entity
}

func traceRelation(stage string, r entity.Relation) {
	fmt.Printf("trace: %s origin=%s target=%s thread=%s impulse=%q\n", stage, r.Origin, r.ID, r.ThreadID, r.Impulse)
}

func (c *ContinueThread) Relate(r entity.Relation) entity.Entity {
	traceRelation("ingress", r)
	return c.relate(r)
}

func (c *ContinueThread) relate(r entity.Relation) entity.Entity {
	name, _ := meaning.Extract(r.Impulse, "~with", "continue-thread")
	message, _ := meaning.Extract(r.Impulse, "~say", "continue-thread", "|")

	target, ok := c.EntityMap[name]
	if !ok {
		fmt.Println("debug: target not found:", name)
		return c
	}
	b, ok := target.(inferrable)
	if !ok {
		fmt.Println("debug: target not inferrable:", name)
		return c
	}

	t, ok := c.EntityMap[r.ThreadID].(Thread)
	if !ok {
		fmt.Println("debug: thread not found:", r.ThreadID)
		return c
	}

	// Only append to the thread if the relation actually carries a message.
	if message != "" {
		t = t.Append(r.Origin, name, r)
		c.EntityMap[r.ThreadID] = t
	}

	// Build present: target's self state + current exchange + other exchanges summary + pulled refs + sender + message
	exchangeLines := t.ExchangeBetween(name, r.Origin)
	threadContext := "\nthread " + t.id + " (" + t.About + "):\n"
	if exchangeLines != "" {
		threadContext += "current exchange with " + r.Origin + ":\n" + exchangeLines
	}
	currentStateContext := ""
	if other := t.OtherExchangesFor(name, r.Origin); other != "" {
		currentStateContext = "\nyour other exchanges in this thread (use ~ref <peer>:<range> to pull):\n" + other
	}
	// Pull referenced entries (from the sender's own exchanges) into the target's present.
	pulledContext := ""
	if ref, _ := meaning.Extract(r.Impulse, "~ref", "continue-thread"); ref != "" {
		entries := t.ResolveRef(r.Origin, ref)
		if len(entries) > 0 {
			pulledContext = "\npulled context (" + ref + "):\n"
			for i, rel := range entries {
				msg, err := meaning.Extract(rel.Impulse, "~say", "exchange", "|")
				if err != nil {
					msg = rel.Impulse
				}
				pulledContext += fmt.Sprintf("  [%d] %s: %s\n", i, rel.Origin, msg)
			}
		}
	}
	senderContext := ""
	messageLine := ""
	if message != "" {
		senderContext = "\nsender: " + r.Origin
		messageLine = "\nmessage from " + r.Origin + ": " + message
	}
	present := b.DerivePresent(r) + threadContext + currentStateContext + pulledContext + senderContext + messageLine

	// Call target's medium to produce a response
	m := b.Medium()
	if m == nil {
		fmt.Println("debug: target has no medium:", name)
		return c
	}
	response, err := m(present, r)
	if err != nil {
		fmt.Println("medium error:", err)
		return c
	}
	if response == "" {
		return c
	}

	// Every line of the response must be a valid protocol string. Non-protocol lines are dropped.
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		next, err := entity.Impress(name, r.ThreadID, line)
		if err != nil {
			fmt.Println("debug: dropping non-protocol line from", name, "→", line)
			continue
		}
		// Reject self-references — a being can't target itself.
		if nextWith, _ := meaning.Extract(next.Impulse, "~with", "continue-thread"); nextWith == name {
			fmt.Println("debug: dropping self-reference from", name)
			continue
		}
		traceRelation("dispatch", next)
		nextNode, ok := c.EntityMap[next.ID]
		if !ok {
			fmt.Println("debug: emitted target not found:", next.ID)
			continue
		}
		if next.ID == c.Name() {
			c.relate(next)
		} else {
			nextNode.Relate(next)
		}
	}
	return c
}

func (c *ContinueThread) ID() string   { return "continue-thread" }
func (c *ContinueThread) Name() string { return "continue-thread" }
