package reality

import (
	"crypto/rand"
	"fmt"
	"skyra-v05/src/debug"
	"strings"
)

type NewThread struct {
	id      string
	Beings  map[string]Reality
	Access  map[string]bool
	Threads map[string]*Thread
}

type Thread struct {
	id        string
	CreatedBy string
	Active    bool
	Members   map[string]bool
	Graph     []Edge
	Queue     []*Relation
}

type Edge struct {
	From string
	To   string
}

func (t *NewThread) ID() string { return t.id }

func (t *NewThread) Create(r *Relation) Reality {
	return &NewThread{
		id:      "thread-gate",
		Beings:  make(map[string]Reality),
		Access:  make(map[string]bool),
		Threads: make(map[string]*Thread),
	}
}

func (t *NewThread) Realize(r *Relation) string {
	if r.ThreadID != "" {
		th, ok := t.Threads[r.ThreadID]
		if !ok {
			debug.Log("[thread]: unknown thread id:", r.ThreadID)
			return ""
		}
		debug.Log("[thread]: passing through existing thread", r.ThreadID)

		if r.ID != "" {
			th.Spread(r.Origin, r.ID)
		}

		r.Attach("thread", th.Parse)
		for name, being := range t.Beings {
			r.Realities[name] = being
		}
		return ""
	}

	if !t.Access[r.Origin] {
		debug.Log("[thread]: no thread access for", r.Origin)
		return ""
	}

	th := t.newThread(r.Origin)
	r.ThreadID = th.id
	debug.Log("[thread]: created thread", th.id, "for", r.Origin)

	if r.ID != "" {
		th.Spread(r.Origin, r.ID)
	}

	r.Attach("thread", th.Parse)
	for name, being := range t.Beings {
		r.Realities[name] = being
	}
	return ""
}

func (t *NewThread) newThread(creator string) *Thread {
	b := make([]byte, 8)
	rand.Read(b)
	th := &Thread{
		id:        fmt.Sprintf("%x", b),
		CreatedBy: creator,
		Active:    true,
		Members:   map[string]bool{creator: true},
	}
	t.Threads[th.id] = th
	return th
}

func (th *Thread) Spread(from, to string) {
	th.Members[from] = true
	if to != "" {
		th.Members[to] = true
	}
	th.Graph = append(th.Graph, Edge{From: from, To: to})
}

func (th *Thread) Parse() string {
	var sb strings.Builder
	sb.WriteString("thread " + th.id + "\n")
	sb.WriteString("created by: " + th.CreatedBy + "\n")
	if th.Active {
		sb.WriteString("status: active\n")
	} else {
		sb.WriteString("status: closed\n")
	}
	sb.WriteString("members:\n")
	for member := range th.Members {
		sb.WriteString("  " + member + "\n")
	}
	return sb.String()
}

func (t *NewThread) Route(origin, target, response, threadID string) []*Relation {
	relations := ParseResponse(origin, response)
	debug.Log("[thread]: route — parsed", len(relations), "emissions")

	if len(relations) == 0 && strings.TrimSpace(response) != "" {
		relations = []*Relation{{
			Origin:    origin,
			ID:        target,
			Impulse:   response,
			Parsers:   make(map[string]Parser),
			Realities: make(map[string]Reality),
		}}
		debug.Log("[thread]: route — plain text fallback →", origin, "to", target)
	}

	th, ok := t.Threads[threadID]
	if ok {
		for _, rel := range relations {
			rel.ThreadID = threadID
			th.Spread(origin, rel.ID)
		}
	}
	return relations
}

func (t *NewThread) Parse() string {
	return ""
}
