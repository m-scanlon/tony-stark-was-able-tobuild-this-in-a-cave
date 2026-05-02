package reality

import (
	"crypto/rand"
	"fmt"
	"skyra-v05/src/debug"
	"strings"
)

type NewThread struct {
	id       string
	Beings   map[string]Reality
	Access   map[string]bool
	Threads  map[string]*Thread
	Exchange *Exchange
	Devices  map[string]Reality
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
	for {
		var th *Thread

		if r.ThreadID != "" {
			var ok bool
			th, ok = t.Threads[r.ThreadID]
			if !ok {
				debug.Log("[thread]: unknown thread id:", r.ThreadID)
				return ""
			}
			debug.Log("[thread]: existing thread", r.ThreadID)
		} else {
			if !t.Access[r.Origin] {
				debug.Log("[thread]: no thread access for", r.Origin)
				return ""
			}
			th = t.newThread(r.Origin)
			r.ThreadID = th.id
			debug.Log("[thread]: created thread", th.id, "for", r.Origin)
		}

		if r.ID == "" {
			op, rest := r.Peel()
			if op == "grow" {
				msg := t.Grow(rest)
				debug.Log("[thread]: grow →", msg)
				if user, ok := t.Beings[r.Origin]; ok {
					r.Impulse = msg
					user.Realize(r)
					r.ID = ""
					r.Origin = r.Origin
					r.Parsers = make(map[string]Parser)
					continue
				}
			}
		}

		r.Attach("thread", th.Parse)

		if r.Realities == nil {
			r.Realities = make(map[string]Reality)
		}
		for name, being := range t.Beings {
			r.Realities[name] = being
		}

		debug.Log("[thread]: descending", r.Origin, "→", r.ID, "|", r.Impulse)

		response := t.Exchange.Realize(r)

		if response == "" {
			if errReality, ok := r.Realities["error"]; ok {
				err := errReality.(*Error)
				debug.Log("[thread]: exchange error →", err.Message)
				delete(r.Realities, "error")

				origin := r.Origin
				errMsg := err.Message
				r.Attach("exchange-error", func() string { return errMsg })

				if being, ok := t.Beings[origin]; ok {
					debug.Log("[thread]: routing error back to", origin)
					response = being.Realize(r)
					if response != "" {
						th.Spread(r.Origin, r.ID)
						r.Parsers = make(map[string]Parser)
						continue
					}
				}
			}
			debug.Log("[thread]: empty response")
			return ""
		}

		th.Spread(r.Origin, r.ID)
		debug.Log("[thread]: routing", r.Origin, "→", r.ID)

		r.Parsers = make(map[string]Parser)
	}
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

func (t *NewThread) Grow(impulse string) string {
	name, err := Extract(impulse, "~name", "grow")
	if err != nil {
		return "grow: missing ~name"
	}
	if _, exists := t.Beings[name]; exists {
		return "grow: " + name + " already exists"
	}
	beingType, err := Extract(impulse, "~type", "grow")
	if err != nil {
		return "grow: missing ~type"
	}
	deviceName, err := Extract(impulse, "~device", "grow")
	if err != nil {
		return "grow: missing ~device"
	}

	device, ok := t.Devices[deviceName]
	if !ok {
		return "grow: unknown device " + deviceName
	}

	being := Being{}.Create(&Relation{
		ID:      name,
		Impulse: impulse,
	}).(Being)

	switch beingType {
	case "llm":
		self := &Self{}
		self = self.Create(&Relation{ID: name}).(*Self)
		self.Realities["being"] = being

		think := &Think{
			Operators: map[string]Reality{
				"recall":   &Recall{},
				"remember": &Remember{},
				"skill":    &Skill{},
			},
			LLM: device,
		}

		act := &Act{
			Operators: map[string]Reality{
				"plan": &Plan{},
			},
			LLM: device,
		}

		self.Realities["think"] = think
		self.Realities["act"] = act
		t.Beings[name] = self

	case "user":
		user := &User{}
		user = user.Create(&Relation{ID: name}).(*User)
		user.Realities["being"] = being
		user.Realities["device"] = device
		t.Beings[name] = user
		t.Access[name] = true

	default:
		return "grow: unknown type " + beingType
	}

	for _, peerName := range being.Relationships {
		if peer, ok := t.Beings[peerName]; ok {
			switch p := peer.(type) {
			case *Self:
				if b, ok := p.Realities["being"].(Being); ok {
					b.Relationships = append(b.Relationships, name)
					p.Realities["being"] = b
				}
			case *User:
				if b, ok := p.Realities["being"].(Being); ok {
					b.Relationships = append(b.Relationships, name)
					p.Realities["being"] = b
				}
			}
		}
	}

	debug.Log("[thread]: grew", name, "type:", beingType, "device:", deviceName)
	return "grew " + name
}

func (t *NewThread) Parse() string {
	return ""
}
