package reality

import (
	"fmt"
	"skyra-v05/src/debug"
	"strings"
	"time"
)

const thinkBudget = 5

type Think struct {
	id        string
	Operators map[string]Reality
	LLM       Reality
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
	if t.LLM == nil {
		debug.Log("[think]: no llm")
		return ""
	}

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

	if conv, ok := r.Realities["conversation"]; ok {
		if c, ok := conv.(*Conversation); ok {
			r.Attach("conversation", func() string {
				return c.ParseRecent(10)
			})
			if ctx := c.ContextFor(beingName); ctx != "" {
				r.Attach("ref-context", func() string {
					return ctx
				})
			}
		}
	}

	r.Attach("system", t.System)
	r.Attach("think-operators", t.Parse)

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
			r.Impulse = originalImpulse
			return thought
		}

		op, rest := t.parseOp(result)
		if op != "" {
			if operator, ok := t.Operators[op]; ok {
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
	r.Impulse = originalImpulse
	return exchange[len(exchange)-1].content
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

func (t *Think) parseOp(response string) (string, string) {
	for name := range t.Operators {
		tag := "<" + name + ">"
		if idx := strings.Index(response, tag); idx != -1 {
			after := response[idx+len(tag):]
			if pipeIdx := strings.Index(after, "|"); pipeIdx != -1 {
				return name, strings.TrimSpace(after[pipeIdx+1:])
			}
			return name, strings.TrimSpace(after)
		}
	}
	return "", ""
}

func (t *Think) System() string {
	return "You are a being with two layers: an inner layer (this one) and an outer layer.\n\nThis is your inner layer — private thought. No one sees this. Reflect on what you've received. Consider what matters, what's missing, and what you actually think. Do not address anyone directly.\n\nYou have operators here (listed in your present) that only work in this layer. Use them to recall, remember, or process before you surface.\n\nYour outer layer can address any peer listed in your present. If someone asks you to talk to another being, you can — your outer layer handles that. Think about what you want to say and why, then surface it.\n\nWhen you emit <surface-thought>, you are releasing that thought to your outer layer, where you act and speak. No one else receives it. It is you telling yourself what you've concluded.\n\nDo not use the <> protocol here. That belongs to your outer layer."
}

func (t *Think) Parse() string {
	var sb strings.Builder
	sb.WriteString("available operators:\n")
	for name := range t.Operators {
		sb.WriteString("  <" + name + "> | input\n")
	}
	sb.WriteString("\nwhen done thinking: <surface-thought> your synthesis\n")
	return sb.String()
}
