package reality

import (
	"skyra-v05/src/debug"
	"strings"
)

type Operators struct {
	id       string
	Registry map[string]func() Reality
}

func NewOperators() *Operators {
	return &Operators{
		id:       "operators",
		Registry: make(map[string]func() Reality),
	}
}

func (o *Operators) ID() string { return o.id }

func (o *Operators) Register(verb string, constructor func() Reality) {
	o.Registry[verb] = constructor
}

func (o *Operators) Create(r *Relation) Reality {
	verb, _ := extractVerb(r.Impulse)
	constructor, ok := o.Registry[verb]
	if !ok {
		return o
	}
	o.Registry[verb] = constructor
	return o
}

func (o *Operators) Realize(r *Relation) string {
	debug.Log("[operators]: impulse →", r.Impulse)
	verb, rest := extractVerb(r.Impulse)
	debug.Log("[operators]: verb →", verb, "| rest →", rest)
	constructor, ok := o.Registry[verb]
	if !ok {
		debug.Log("[operators]: verb not found, falling back to default (exchange)")
		constructor, ok = o.Registry["exchange"]
		if !ok {
			debug.Log("[operators]: no exchange registered")
			return ""
		}
		rest = r.Impulse
	}

	r.Attach("operators", o.Parse)

	op := constructor()
	r.Impulse = rest
	return op.Realize(r)
}

func (o *Operators) Parse() string {
	var sb strings.Builder
	sb.WriteString("available operators:\n")
	for verb := range o.Registry {
		sb.WriteString("  " + verb + "\n")
	}
	return sb.String()
}

func extractVerb(impulse string) (string, string) {
	for i, c := range impulse {
		if c == ' ' {
			return impulse[:i], impulse[i+1:]
		}
	}
	return impulse, ""
}
