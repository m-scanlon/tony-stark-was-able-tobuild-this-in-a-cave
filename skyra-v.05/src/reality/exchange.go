package reality

import (
	"fmt"
	"skyra-v05/src/debug"
	"strings"
)

type Exchange struct {
	id        string
	Exchanges map[string]*Conversation
}

type Conversation struct {
	Parties [2]string
	Active  bool
	Entries []Entry
	Sender  string
	Message string
}

type Entry struct {
	From    string
	Content string
}

func (e *Exchange) ID() string { return e.id }

func (e *Exchange) Create(r *Relation) Reality {
	return &Exchange{
		id:        "exchange",
		Exchanges: make(map[string]*Conversation),
	}
}

func exchangeKey(a, b string) string {
	if a < b {
		return a + ":" + b
	}
	return b + ":" + a
}

func (e *Exchange) Realize(r *Relation) string {
	if r.ID == "" {
		target, rest := r.Peel()
		debug.Log("[exchange]: peeled target →", target, "| rest →", rest)
		r.ID = target
		r.Impulse = rest
	} else {
		debug.Log("[exchange]: target already set →", r.ID)
	}

	key := exchangeKey(r.Origin, r.ID)
	conv, ok := e.Exchanges[key]
	if !ok {
		conv = &Conversation{
			Parties: [2]string{r.Origin, r.ID},
			Active:  true,
		}
		e.Exchanges[key] = conv
		debug.Log("[exchange]: new conversation", key)
	} else {
		debug.Log("[exchange]: existing conversation", key, "| entries:", len(conv.Entries))
	}

	conv.Sender = r.Origin
	conv.Message = r.Impulse
	debug.Log("[exchange]: recording entry from", r.Origin, "→", r.Impulse[:min(len(r.Impulse), 40)])
	conv.Entries = append(conv.Entries, Entry{
		From:    r.Origin,
		Content: r.Impulse,
	})
	r.Attach("exchange", conv.Parse)

	being, ok := r.Realities[r.ID]
	if !ok {
		debug.Log("[exchange]: being not found:", r.ID, "| available:", len(r.Realities))
		for name := range r.Realities {
			debug.Log("[exchange]:   -", name)
		}
		return ""
	}
	debug.Log("[exchange]: routing to being:", r.ID)
	return being.Realize(r)
}

func (c *Conversation) Parse() string {
	var sb strings.Builder
	sb.WriteString("exchange: " + c.Parties[0] + " ↔ " + c.Parties[1] + "\n")
	if c.Active {
		sb.WriteString("status: active\n")
	} else {
		sb.WriteString("status: closed\n")
	}
	for i, entry := range c.Entries {
		sb.WriteString(fmt.Sprintf("  [%d] %s: %s\n", i, entry.From, entry.Content))
	}
	if c.Sender != "" {
		sb.WriteString("\nsender: " + c.Sender + "\n")
		sb.WriteString("message from " + c.Sender + ": " + c.Message + "\n")
	}
	return sb.String()
}

func (e *Exchange) Parse() string {
	return ""
}
