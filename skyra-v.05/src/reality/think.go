package reality

import (
	"fmt"
	"skyra-v05/src/debug"
	"strings"
	"time"
)

const thinkBudget = 5
const thoughtHistoryMax = 10

type Think struct {
	id        string
	Operators map[string]Reality
	OuterOps  []string
	LLM       Reality
	History   []ThoughtSection
}

type ThoughtSection struct {
	Peer      string
	Thought   string
	Timestamp time.Time
}

type thinkEntry struct {
	timestamp time.Time
	speaker   string
	content   string
}

func (t *Think) ID() string { return t.id }

func (t *Think) Create(r *Relation) Reality {
	return &Think{
		id:        "think",
		Operators: make(map[string]Reality),
	}
}

func (t *Think) Realize(r *Relation) string {
	if r.Collecting {
		snap := ThinkSnapshot{
			Budget:    thinkBudget,
			Operators: []string{},
			History:   []ThoughtSnapshot{},
		}
		for name := range t.Operators {
			snap.Operators = append(snap.Operators, name)
		}
		for _, h := range t.History {
			snap.History = append(snap.History, ThoughtSnapshot{
				Peer: h.Peer, Thought: h.Thought, Ts: h.Timestamp.UnixMilli(),
			})
		}
		r.Export("think", snap)

		node := RealityNode{ID: "think", Type: "Think", Children: []RealityNode{}}
		for name := range t.Operators {
			node.Children = append(node.Children, RealityNode{
				ID: name, Type: capitalizeType(name), Children: []RealityNode{},
			})
		}
		r.Export("node:think", node)
		return ""
	}

	if t.LLM == nil {
		debug.Log("[think]: no llm")
		return ""
	}

	ops := t.collectOps(r)

	var beingName string
	if being, ok := r.Realities["being"]; ok {
		if b, ok := being.(Being); ok {
			beingName = b.Name()
			r.Attach("being", b.ParseInner)
		}
	}
	if beingName == "" {
		beingName = "self"
	}

	log := func(args ...any) { debug.Being(beingName, "inner", args...) }
	r.Log = log

	r.Attach("system", t.System)
	outerOps := t.OuterOps
	r.Attach("think-operators", func() string { return renderOpsWithOuter(ops, outerOps) })

	if len(t.History) > 0 {
		history := t.History
		r.Attach("thought-history", func() string {
			return renderThoughtHistory(history)
		})
	}

	peer := r.Origin

	originalImpulse := r.Impulse
	var exchange []thinkEntry

	for i := 0; i < thinkBudget; i++ {
		log("[think]: pass", i)

		r.Impulse = originalImpulse

		remaining := thinkBudget - i
		r.Attach("think-time", func() string {
			return timePressure(remaining)
		})

		if remaining == 1 {
			delete(r.Parsers, "think-operators")
		}

		ex := exchange
		r.Attach("think-exchange", func() string {
			return renderThinkExchange(ex)
		})

		result := t.LLM.Realize(r)
		log("[think]: llm returned →", result)

		exchange = append(exchange, thinkEntry{
			timestamp: time.Now(),
			speaker:   beingName,
			content:   result,
		})

		thought, done := parseThink(result)
		if done {
			log("[think]: done after", i+1, "passes")
			t.recordThought(peer, stripSurface(thought))
			r.Impulse = originalImpulse
			return thought
		}

		if t.isOuterOp(result) != "" {
			blocked := t.isOuterOp(result)
			log("[think]: blocked outer operator", blocked)
			exchange = append(exchange, thinkEntry{
				timestamp: time.Now(),
				speaker:   "system",
				content:   fmt.Sprintf("<%s> belongs to your outer layer. you cannot call it here. use your inner operators or surface your thought.", blocked),
			})
			continue
		}

		op, rest := parseOp(result, ops)
		if op != "" {
			if operator, ok := ops[op]; ok {
				log("[think]: firing operator", op)
				r.Impulse = rest
				opResult := operator.Realize(r)
				if opResult != "" {
					exchange = append(exchange, thinkEntry{
						timestamp: time.Now(),
						speaker:   op,
						content:   opResult,
					})
				}
			} else {
				log("[think]: unknown operator", op)
			}
			continue
		}
	}

	log("[think]: budget exhausted")
	last := exchange[len(exchange)-1].content
	t.recordThought(peer, last)
	r.Impulse = originalImpulse
	return last
}

func (t *Think) recordThought(peer, thought string) {
	t.History = append(t.History, ThoughtSection{
		Peer:      peer,
		Thought:   strings.TrimSpace(thought),
		Timestamp: time.Now(),
	})
	if len(t.History) > thoughtHistoryMax {
		t.History = t.History[len(t.History)-thoughtHistoryMax:]
	}
}

func renderThoughtHistory(sections []ThoughtSection) string {
	var sb strings.Builder
	sb.WriteString("your recent thoughts (across exchanges):\n")
	for _, s := range sections {
		sb.WriteString(fmt.Sprintf("[%s] (with %s): %s\n\n", s.Timestamp.Format("15:04:05"), s.Peer, s.Thought))
	}
	return sb.String()
}

func renderThinkExchange(entries []thinkEntry) string {
	if len(entries) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("your previous thoughts in this thinking session:\n")
	for _, e := range entries {
		sb.WriteString(e.timestamp.Format("15:04:05") + "\n")
		sb.WriteString(e.speaker + ": " + e.content + "\n\n")
	}
	return sb.String()
}

func timePressure(remaining int) string {
	switch remaining {
	case 1:
		return fmt.Sprintf("time remaining: %d. You must emit <surface-thought> now.\n", remaining)
	case 2:
		return fmt.Sprintf("time remaining: %d. Wrap up your thinking.\n", remaining)
	default:
		return fmt.Sprintf("time remaining: %d.\n", remaining)
	}
}

func parseThink(response string) (string, bool) {
	return response, strings.Contains(response, "<surface-thought>")
}

func parseOp(response string, ops map[string]Reality) (string, string) {
	for name := range ops {
		openTag := "<" + name + ">"
		closeTag := "</" + name + ">"
		if idx := strings.Index(response, openTag); idx != -1 {
			after := response[idx+len(openTag):]
			if end := strings.Index(after, closeTag); end != -1 {
				return name, strings.TrimSpace(after[:end])
			}
			if pipeIdx := strings.Index(after, "|"); pipeIdx != -1 {
				return name, strings.TrimSpace(after[pipeIdx+1:])
			}
			return name, strings.TrimSpace(after)
		}
	}
	return "", ""
}

func (t *Think) collectOps(r *Relation) map[string]Reality {
	ops := make(map[string]Reality)
	for name, op := range t.Operators {
		ops[name] = op
	}
	if r.Realities != nil {
		for key, op := range r.Realities {
			if strings.HasPrefix(key, "think:") {
				name := strings.TrimPrefix(key, "think:")
				ops[name] = op
			}
		}
	}
	return ops
}

func renderOps(ops map[string]Reality) string {
	var sb strings.Builder
	sb.WriteString("available operators:\n")
	for name := range ops {
		sb.WriteString("  <" + name + ">input</" + name + ">\n")
	}
	sb.WriteString("\nwhen done thinking: <surface-thought> your synthesis\n")
	return sb.String()
}

func renderOpsWithOuter(ops map[string]Reality, outerOps []string) string {
	var sb strings.Builder
	sb.WriteString("available operators:\n")
	for name := range ops {
		sb.WriteString("  <" + name + ">input</" + name + ">\n")
	}
	if len(outerOps) > 0 {
		sb.WriteString("\nouter layer operators (you cannot call these here):\n")
		for _, name := range outerOps {
			sb.WriteString("  " + name + "\n")
		}
	}
	sb.WriteString("\nwhen done thinking: <surface-thought> your synthesis\n")
	return sb.String()
}

func (t *Think) System() string {
	return "You are a being with two layers: an inner layer (this one) and an outer layer.\n\nThis is your inner layer — private thought. No one sees this. Reflect on what you've received. Consider what matters, what's missing, and what you actually think. Do not address anyone directly.\n\nYou have operators here (listed in your present) that only work in this layer. Use them to recall, remember, or process before you surface.\n\nYour outer layer can address any peer listed in your present. If someone asks you to talk to another being, you can — your outer layer handles that. Think about what you want to say and why, then surface it.\n\nIMPORTANT: Emit exactly one protocol per response — one operator OR one <surface-thought>. Never both. If you call an operator, wait for the result before doing anything else. You have multiple passes to think.\n\nWhen you emit <surface-thought>, you are releasing that thought to your outer layer, where you act and speak. No one else receives it. It is you telling yourself what you've concluded.\n\nDo not use the <> protocol here. That belongs to your outer layer."
}

func (t *Think) isOuterOp(response string) string {
	for _, name := range t.OuterOps {
		tag := "<" + name + ">"
		if strings.Contains(response, tag) {
			return name
		}
	}
	return ""
}

func (t *Think) Parse() string {
	return renderOps(t.Operators)
}
