package reality

import (
	"strings"

	"skyra-v05/src/debug"
)

type Act struct {
	id        string
	Operators map[string]Reality
	LLM       Reality
}

func (a *Act) ID() string { return a.id }

func (a *Act) Create(r *Relation) Reality {
	return &Act{
		id:        "act",
		Operators: make(map[string]Reality),
	}
}

func (a *Act) Realize(r *Relation) string {
	r.Attach("system", a.System)

	var beingName string
	if being, ok := r.Realities["being"]; ok {
		if b, ok := being.(Being); ok {
			beingName = b.Name()
			r.Attach("being", b.Parse)
		}
	}
	if beingName == "" {
		beingName = "self"
	}

	log := func(args ...any) { debug.Being(beingName, "outer", args...) }
	r.Log = log

	if conv, ok := r.Realities["conversation"]; ok {
		if c, ok := conv.(*Conversation); ok {
			if ctx := c.ContextFor(beingName); ctx != "" {
				r.Attach("ref-context", func() string {
					return ctx
				})
			}
		}
	}

	if innerParser, ok := r.Parsers["inner"]; ok {
		raw := innerParser()
		synthesis := stripSurface(raw)
		log("[act]: synthesis →", synthesis)
		r.Attach("inner", func() string { return "your inner thought: " + synthesis + "\n" })
	}

	if a.LLM == nil {
		log("[act]: no llm")
		return ""
	}

	for attempt := 0; attempt < 3; attempt++ {
		log("[act]: calling llm, attempt", attempt)
		result := a.LLM.Realize(r)
		log("[act]: response →", result)

		relations := ParseResponse("", result)
		selfRoute := false
		for _, rel := range relations {
			if rel.ID == beingName {
				selfRoute = true
				break
			}
		}

		if len(relations) == 0 {
			warning := "WARNING: your response did not follow the protocol. You must use: <> target | message. Try again.\n"
			log("[act]: protocol violation, retrying")
			r.Attach("act-warning", func() string { return warning })
			r.Impulse = result
			continue
		}

		if selfRoute {
			warning := "WARNING: you addressed yourself. You cannot route messages to yourself. Address a peer instead. Try again.\n"
			log("[act]: self-route detected, retrying")
			r.Attach("act-warning", func() string { return warning })
			r.Impulse = result
			continue
		}

		delete(r.Parsers, "act-warning")
		r.Origin = ""
		return result
	}

	log("[act]: exhausted retries")
	r.Origin = ""
	return r.Impulse
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

func (a *Act) System() string {
	return "You are a being with two layers: an inner layer and an outer layer (this one).\n\nThis is your outer layer — where you act and speak. Your inner layer has already processed and given you a thought (shown in your present). Now you respond.\n\nPROTOCOL\nEvery response is exactly one message in this format:\n  <> target | message\n\nExamples:\n  <> michael | hello, what's on your mind?\n  <> builder | can you check the deployment?\n\nOnly one <> line per response. If you're replying to whoever sent you a message, use their name as the target.\n\nWhen you leave a conversation to talk to someone else, carry context with ~ref:\n  <> louise | hey, wanted to talk ~ref michael:0-3\nThis brings entries 0-3 from your current exchange as your private context in the new conversation.\n\nIMPORTANT: To talk to a peer, emit a message to them directly. Do NOT say \"I will go talk to them\" — that doesn't do anything. Actually address them.\n\nDo not use operators like <recall> or <remember> here. Those belong to your inner layer.\n\nNever start your response with your own name. No asterisks, no roleplay, no action narration."
}

func (a *Act) Parse() string {
	return ""
}
