package world

import (
	"fmt"
	"strings"

	"skyra-v04/src/primitives/being"
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
	"skyra-v04/src/primitives/thread"
)

type World struct {
	EntityMap map[string]entity.Entity
	threads   map[string]*thread.Thread
	id        string
	name      string
}

func New() *World {
	return &World{
		EntityMap: make(map[string]entity.Entity),
		threads:   make(map[string]*thread.Thread),
	}
}

func (w *World) ID() string   { return w.id }
func (w *World) Name() string { return w.name }

func (w *World) Relate(r entity.Relation) entity.Entity {
	switch r.ID {
	case "grow":
		w.grow(r)
		return w

	default:
		if r.ThreadID == "" {
			threadID := thread.NewThreadID()
			r.ThreadID = threadID
			about, _ := meaning.Extract(r.Impulse, "~about", "thread")
			because, _ := meaning.Extract(r.Impulse, "~because", "thread")
			r.Impulse = meaning.Strip(r.Impulse, "~about")
			r.Impulse = meaning.Strip(r.Impulse, "~because")
			t, _ := thread.Thread{}.Relate(entity.Relation{
				ThreadID: threadID,
				Impulse:  fmt.Sprintf("~about %s ~because %s", about, because),
			}).(thread.Thread)
			w.threads[threadID] = &t
		}
		w.DerivePresent(r)
		return w
	}
}

func (w *World) DerivePresent(r entity.Relation) string {
	traceRelation("ingress", r)

	name := r.ID
	target, ok := w.EntityMap[name]
	if !ok {
		fmt.Println("debug: target not found:", name)
		return ""
	}
	b, ok := target.(being.IBeing)
	if !ok {
		fmt.Println("debug: target not a being:", name)
		return ""
	}

	t, ok := w.threads[r.ThreadID]
	if !ok {
		fmt.Println("debug: thread not found:", r.ThreadID)
		return ""
	}

	ref, _ := meaning.Extract(r.Impulse, "~ref", "route")
	cleanImpulse := meaning.Strip(r.Impulse, "~ref")
	// Strip target name from the message — the model sometimes echoes
	// the recipient's name at the start (e.g. "<michael> hi" or "michael, hi").
	cleanImpulse = stripTargetName(cleanImpulse, name)
	storageRel := r
	storageRel.Impulse = cleanImpulse

	if cleanImpulse != "" {
		*t = t.Append(r.Origin, name, storageRel)
	}

	present := w.derivePresent(b, t, r, cleanImpulse, ref)
	logPresent(name, present)
	logThreadState(name, "before-medium", *t)

	m := b.Medium()
	if m == nil {
		fmt.Println("debug: target has no medium:", name)
		return ""
	}
	response, err := m(present, r)
	if err != nil {
		fmt.Println("medium error:", err)
		return ""
	}
	if response == "" {
		return ""
	}

	valid, formatErrs := w.parseResponse(response, name, r.ThreadID, r.Origin, present)
	if len(formatErrs) > 0 {
		fmt.Println("debug:", name, "had", len(formatErrs), "invalid lines — retrying with feedback")
		feedback := present + "\n\n=== RETRY FEEDBACK (this is internal guidance to you, NOT a message to send, do NOT include any of this text in your response) ===\nYour previous response produced these invalid lines, which were dropped:\n"
		for _, e := range formatErrs {
			feedback += "  - " + e + "\n"
		}
		feedback += "Please retry. Just respond naturally and end with <>. Your message goes to whoever you're currently in exchange with. To address a different peer, start with their name.\n=== END RETRY FEEDBACK ===\n"
		if retry, rerr := m(feedback, r); rerr == nil && retry != "" {
			valid, formatErrs = w.parseResponse(retry, name, r.ThreadID, r.Origin, feedback)
			for _, e := range formatErrs {
				fmt.Println("debug: dropping after retry —", e)
			}
		}
	}

	for _, next := range valid {
		nextTarget := next.ID

		if nextTarget == name {
			reason := "self-reference — a being cannot target itself"
			fmt.Println("debug: dropping self-reference from", name)
			logDrop(name, nextTarget+" "+next.Impulse, reason, present)
			continue
		}

		autoClose := false
		if returnTo := t.FindReturnTarget(name); returnTo != "" && nextTarget == returnTo {
			fmt.Println("debug: auto-close from", name, "→ returning to", returnTo)
			*t = t.CloseExchange(name, r.Origin)
			autoClose = true
		}

		hasRef, _ := meaning.Extract(next.Impulse, "~ref", "dispatch")
		if !autoClose && hasRef != "" && nextTarget != r.Origin {
			fmt.Println("debug: ~ref return from", name, "→ closing exchange with", r.Origin, "→ routing to", nextTarget)
			*t = t.CloseExchange(name, r.Origin)
			autoClose = true
		}

		if !autoClose && nextTarget != r.Origin {
			if ex, exists := t.ExchangeWith(name, nextTarget); exists && ex.Active && ex.Parent != name {
				reason := fmt.Sprintf("blocked message to %s — that peer is your parent; address them to return", nextTarget)
				fmt.Println("debug: blocked from", name, "to", nextTarget)
				logDrop(name, nextTarget+" "+next.Impulse, reason, present)
				continue
			}
		}

		traceRelation("dispatch", next)

		if _, ok := w.EntityMap[next.ID]; !ok {
			reason := fmt.Sprintf("target %q not found in world", next.ID)
			fmt.Println("debug: target not found:", next.ID)
			logDrop(name, next.ID+" "+next.Impulse, reason, present)
			continue
		}

		w.DerivePresent(next)
	}
	return ""
}

func (w *World) derivePresent(b being.IBeing, t *thread.Thread, r entity.Relation, cleanImpulse, ref string) string {
	name := b.Name()

	threadContext := "\nthread " + t.ID() + " (" + t.About + "):\n"

	activeExchanges := ""
	if ae := t.ActiveExchangesFor(name, r.Origin); ae != "" {
		activeExchanges = "\nactive exchanges:\n" + ae
	}

	// Render ~ref context inside the active exchanges section.
	// This is the being's own pulled context — personal to them,
	// attached to the exchange it came from.
	if ref != "" {
		entries := t.ResolveRef(r.Origin, ref)
		if len(entries) > 0 {
			activeExchanges += "  context you pulled (~ref " + ref + "):\n"
			for i, rel := range entries {
				activeExchanges += fmt.Sprintf("    [%d] %s: %s\n", i, rel.Origin, rel.Impulse)
			}
		}
	}

	exchangeLines := t.ExchangeBetween(name, r.Origin)
	currentExchange := ""
	if exchangeLines != "" {
		currentExchange = "\nexchange with " + r.Origin + ":\n" + exchangeLines
	}

	peersContext := buildPeers(name, b.Relationships())
	peersContext += "\n<> <peer> <expression> ~ref <current-peer>:START-END\n"

	senderContext := ""
	messageLine := ""
	if cleanImpulse != "" {
		senderContext = "\nsender: " + r.Origin
		messageLine = "\nmessage from " + r.Origin + ": " + cleanImpulse
	}

	return b.DerivePresent(r) + threadContext + activeExchanges + currentExchange + peersContext + senderContext + messageLine
}

func (w *World) parseResponse(response, beingName, threadID, currentPeer, present string) ([]entity.Relation, []string) {
	var valid []entity.Relation
	var errs []string

	// Split on <> delimiter — each message starts with <>
	parts := strings.Split(response, "<>")
	var blocks []string
	for _, p := range parts {
		b := strings.TrimSpace(p)
		if b != "" {
			blocks = append(blocks, b)
		}
	}

	for _, block := range blocks {
		if block == "" {
			continue
		}

		// Strip the being's own name if it prefixed its response with it.
		tokens := strings.Fields(block)
		firstWord := strings.ToLower(strings.TrimRight(tokens[0], ",:;."))
		if firstWord == beingName && len(tokens) > 1 {
			block = strings.TrimSpace(strings.Join(tokens[1:], " "))
			tokens = strings.Fields(block)
			if len(tokens) == 0 {
				continue
			}
			firstWord = strings.ToLower(strings.TrimRight(tokens[0], ",:;."))
		}

		// Check if the first word is a known peer. If not, this is an
		// untargeted message — route it to the current exchange partner.
		var rel entity.Relation
		if _, isPeer := w.EntityMap[firstWord]; isPeer && firstWord != beingName {
			var err error
			rel, err = entity.Impress(beingName, threadID, block)
			if err != nil {
				reason := fmt.Sprintf("impress error: %v", err)
				errs = append(errs, fmt.Sprintf("%q — %v", block, err))
				logDrop(beingName, block, reason, present)
				continue
			}
		} else {
			if currentPeer == "" {
				reason := "no current peer to route untargeted message to"
				errs = append(errs, fmt.Sprintf("%q — %s", block, reason))
				logDrop(beingName, block, reason, present)
				continue
			}
			rel = entity.Relation{
				ID:       currentPeer,
				Origin:   beingName,
				ThreadID: threadID,
				Impulse:  block,
			}
		}

		if rel.ID == beingName {
			reason := "targets yourself; a being cannot route to itself"
			errs = append(errs, fmt.Sprintf("%q — %s", block, reason))
			logDrop(beingName, block, reason, present)
			continue
		}

		if currentPeer != "" && rel.ID != currentPeer {
			if _, err := meaning.Extract(rel.Impulse, "~ref", "parse"); err != nil {
				reason := fmt.Sprintf("messaging %s outside current exchange; ~ref %s:START-END required", rel.ID, currentPeer)
				errs = append(errs, fmt.Sprintf("%q — %s", block, reason))
				logDrop(beingName, block, reason, present)
				continue
			}
		}

		valid = append(valid, rel)
	}
	return valid, errs
}

func stripTargetName(impulse, target string) string {
	impulse = strings.TrimSpace(impulse)
	if impulse == "" {
		return ""
	}
	lower := strings.ToLower(impulse)

	// <name> prefix
	tag := "<" + target + ">"
	if strings.HasPrefix(lower, tag) {
		return strings.TrimSpace(impulse[len(tag):])
	}

	// name followed by punctuation or space
	if strings.HasPrefix(lower, target) {
		rest := impulse[len(target):]
		if rest == "" {
			return ""
		}
		if rest[0] == ',' || rest[0] == ':' || rest[0] == ';' || rest[0] == '.' || rest[0] == ' ' {
			return strings.TrimSpace(strings.TrimLeft(rest, ",:;. "))
		}
	}

	return impulse
}

func buildPeers(beingID string, relationships []string) string {
	var sb strings.Builder
	sb.WriteString("\npeers you can address:\n")
	for _, peer := range relationships {
		if peer == beingID {
			continue
		}
		sb.WriteString("  " + peer + "\n")
	}
	return sb.String()
}
