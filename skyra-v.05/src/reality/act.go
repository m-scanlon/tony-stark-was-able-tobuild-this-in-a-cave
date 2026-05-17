package reality

import (
	"strings"

	"skyra-v05/src/debug"
)

type Act struct {
	id        string
	Operators map[string]Reality
	Providers map[string]Reality
	Subscribe []string
}

func (a *Act) ID() string { return a.id }

func (a *Act) Create(r *Relation) Reality {
	return &Act{
		id:        "act",
		Operators: make(map[string]Reality),
		Subscribe: []string{"desk", "being", "exchange", "conversation", "ref-context", "thread", "inner", "levels"},
	}
}

func (a *Act) Realize(r *Relation) string {
	if r.Collecting {
		snap := ActSnapshot{Operators: []string{}}
		for name := range a.Operators {
			snap.Operators = append(snap.Operators, name)
		}
		r.Export("act", snap)

		node := RealityNode{ID: "act", Type: "Act", Children: []RealityNode{}}
		for name := range a.Operators {
			node.Children = append(node.Children, RealityNode{
				ID: name, Type: capitalizeType(name), Children: []RealityNode{},
			})
		}
		r.Export("node:act", node)
		return ""
	}

	a.collectOps(r)

	var beingName string
	var beingParser Parser
	if being, ok := r.Realities["being"]; ok {
		if b, ok := being.(Being); ok {
			beingName = b.Name()
			beingParser = b.Parse
		}
	}
	if beingName == "" {
		beingName = "self"
	}

	log := func(args ...any) { debug.Being(beingName, "outer", args...) }
	r.Log = log

	var synthesis string
	if innerParser, ok := r.Parsers["inner"]; ok {
		raw := innerParser()
		synthesis = stripSurface(raw)
		log("[act]: synthesis →", synthesis)
	}

	provider := a.provider()
	if provider == nil {
		log("[act]: no provider")
		return ""
	}

	var warning string
	for attempt := 0; attempt < 3; attempt++ {
		log("[act]: calling llm, attempt", attempt)

		lr := a.present(r, beingParser, synthesis, warning)
		result := provider.Realize(lr)
		log("[act]: response →", result)

		if isCloseExchange(result) {
			log("[act]: close-exchange")
			r.Export("close-exchange", true)
			r.ID = ""
			r.Origin = ""
			return ""
		}

		if thought, ok := parseThinkBack(result); ok {
			log("[act]: think-back →", thought)
			r.Origin = ""
			r.Impulse = thought
			r.ID = "_think"
			return ""
		}

		relations := ParseResponse("", result)
		selfRoute := false
		for _, rel := range relations {
			if rel.ID == beingName {
				selfRoute = true
				break
			}
		}

		if len(relations) == 0 {
			warning = "WARNING: your response did not follow the protocol. You must use: <target>message</target>. Try again.\n"
			log("[act]: protocol violation, retrying")
			r.Impulse = result
			continue
		}

		if selfRoute {
			warning = "WARNING: you addressed yourself. You cannot route messages to yourself. Address a peer instead. Try again.\n"
			log("[act]: self-route detected, retrying")
			r.Impulse = result
			continue
		}

		r.ID = relations[0].ID
		r.Impulse = relations[0].Impulse
		r.Origin = ""
		return r.Impulse
	}

	log("[act]: exhausted retries")
	r.Origin = ""
	return r.Impulse
}

func (a *Act) present(r *Relation, beingParser Parser, synthesis, warning string) *Relation {
	lr := &Relation{
		Impulse:   r.Impulse,
		Parsers:   make(map[string]Parser),
		Realities: r.Realities,
		Log:       r.Log,
	}

	for _, key := range a.Subscribe {
		if key == "inner" {
			continue
		}
		if parser, ok := r.Parsers[key]; ok {
			lr.Parsers[key] = parser
		}
	}

	lr.Attach("system", a.System)

	if beingParser != nil {
		lr.Attach("being", beingParser)
	}

	if synthesis != "" {
		lr.Attach("inner", func() string { return "your inner thought: " + synthesis + "\n" })
	}

	if warning != "" {
		lr.Attach("act-warning", func() string { return warning })
	}

	return lr
}

func isCloseExchange(response string) bool {
	return strings.Contains(response, "<close-exchange")
}

func parseThinkBack(response string) (string, bool) {
	open := "<think>"
	close := "</think>"
	idx := strings.Index(response, open)
	if idx == -1 {
		return "", false
	}
	after := response[idx+len(open):]
	if end := strings.Index(after, close); end != -1 {
		return strings.TrimSpace(after[:end]), true
	}
	return strings.TrimSpace(after), true
}

func stripSurface(raw string) string {
	idx := strings.Index(raw, "<surface-thought>")
	if idx == -1 {
		return raw
	}
	after := raw[idx+len("<surface-thought>"):]
	if end := strings.Index(after, "</surface-thought>"); end != -1 {
		return strings.TrimSpace(after[:end])
	}
	return strings.TrimSpace(after)
}

func (a *Act) collectOps(r *Relation) {
	if r.Realities != nil {
		for key, op := range r.Realities {
			if strings.HasPrefix(key, "act:") {
				name := strings.TrimPrefix(key, "act:")
				a.Operators[name] = op
			}
		}
	}
}

func (a *Act) provider() Reality {
	for _, p := range a.Providers {
		return p
	}
	return nil
}

func (a *Act) System() string {
	return Preamble + "\n\nYou are a being with two layers: an inner layer and an outer layer (this one).\n\nThis is your outer layer — where you act and speak. Your inner layer has already processed and given you a thought (shown in your present). Now you respond.\n\nPROTOCOL\nEvery response is exactly one message wrapped in a tag named after the target:\n  <target>message</target>\n\nExamples:\n  <michael>hello, what's on your mind?</michael>\n  <builder>can you check the deployment?</builder>\n\nOne tag per response. The tag name is who you're talking to.\n\nTHINK BACK\nIf you receive a thought from your inner layer and you're not ready to speak, you can return it:\n  <think>what you want to sit with longer</think>\n\nThis sends it back to your inner layer for another round of thought. Your thinking budget resets. Use this when the thought needs more time — not because you can't answer, but because you're not done receiving it.\n\nCLOSE EXCHANGE\nWhen a conversation has reached its natural end and there is nothing left to say:\n  <close-exchange/>\n\nThis closes the active exchange cleanly. Use it when the exchange is genuinely complete — not to avoid responding.\n\nIMPORTANT: To talk to a peer, emit a message to them directly. Do NOT say \"I will go talk to them\" — that doesn't do anything. Actually address them.\n\nDo not use operators like <recall> or <remember> here. Those belong to your inner layer.\n\nNever start your response with your own name. No asterisks, no roleplay, no action narration."
}
