package reality

import (
	"crypto/rand"
	"fmt"
	"skyra-v05/src/debug"
	"strings"
)

type NewThread struct {
	id        string
	Beings    map[string]Reality
	Access    map[string]bool
	Threads   map[string]*Thread
	Exchange  *Exchange
	Levels    *Levels
	Devices   map[string]Reality
	ThinkOps  map[string]Reality
	ActOps    map[string]Reality
	OnResolve func()
}

type Thread struct {
	id        string
	CreatedBy string
	Active    bool
	Members   map[string]bool
	Graph     []Edge
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
	if r.Collecting {
		t.Exchange.Realize(r)
		if t.Levels != nil {
			t.Levels.Realize(r)
		}

		root := RealityNode{ID: "newthread", Type: "NewThread", Children: []RealityNode{}}
		if node, ok := r.Exports["node:exchange"]; ok {
			root.Children = append(root.Children, node.(RealityNode))
			delete(r.Exports, "node:exchange")
		}

		for name, device := range t.Devices {
			device.Realize(r)
			if node, ok := r.Exports["node:"+name]; ok {
				root.Children = append(root.Children, node.(RealityNode))
				delete(r.Exports, "node:"+name)
			}
		}

		for name, being := range t.Beings {
			being.Realize(r)
			if node, ok := r.Exports["node:"+name]; ok {
				root.Children = append(root.Children, node.(RealityNode))
				delete(r.Exports, "node:"+name)
			}
		}

		for _, th := range t.Threads {
			snap := ThreadSnapshot{
				ID:        th.id,
				CreatedBy: th.CreatedBy,
				Active:    th.Active,
				Members:   []string{},
				Edges:     []EdgeSnapshot{},
			}
			for member := range th.Members {
				snap.Members = append(snap.Members, member)
			}
			for _, edge := range th.Graph {
				snap.Edges = append(snap.Edges, EdgeSnapshot{From: edge.From, To: edge.To})
			}
			r.Export("thread:"+th.id, snap)
		}

		r.Export("node:root", root)
		return ""
	}

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
			if op == "being" || op == "grow" {
				msg := t.Grow(rest)
				debug.Log("[thread]: grow →", msg)
				if user, ok := t.Beings[r.Origin]; ok {
					r.Impulse = msg
					user.Realize(r)
					r.ID = ""
					r.Parsers = make(map[string]Parser)
					continue
				}
			}
			if op == "accept" || op == "reject" {
				msg := t.AcceptReject(op, rest, r.Origin)
				debug.Log("[thread]:", op, "→", msg)
				if user, ok := t.Beings[r.Origin]; ok {
					r.Impulse = msg
					user.Realize(r)
					r.ID = ""
					r.Parsers = make(map[string]Parser)
					continue
				}
			}
		}

		r.Attach("thread", th.Parse)

		if t.Levels != nil {
			levels := t.Levels
			self := r.ID
			peer := r.Origin
			r.Attach("levels", func() string { return levels.ParseFor(self, peer) })
		}

		if r.Realities == nil {
			r.Realities = make(map[string]Reality)
		}
		for name, being := range t.Beings {
			r.Realities[name] = being
		}
		for name, op := range t.ThinkOps {
			r.Realities["think:"+name] = op
		}
		for name, op := range t.ActOps {
			r.Realities["act:"+name] = op
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
						if t.OnResolve != nil {
							t.OnResolve()
						}
						r.Parsers = make(map[string]Parser)
						continue
					}
				}
			}
			debug.Log("[thread]: empty response")
			return ""
		}

		th.Spread(r.Origin, r.ID)
		if t.OnResolve != nil {
			t.OnResolve()
		}
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

func (t *NewThread) AcceptReject(op, rest, origin string) string {
	tokens := strings.Fields(rest)
	if len(tokens) < 2 {
		return op + ": usage: " + op + " <being> <task name>"
	}
	beingName := tokens[0]
	taskName := strings.Join(tokens[1:], " ")

	being, ok := t.Beings[beingName]
	if !ok {
		return op + ": being " + beingName + " not found"
	}

	self, ok := being.(*Self)
	if !ok {
		return op + ": " + beingName + " has no desk"
	}

	desk, ok := self.Realities["desk"].(*Desk)
	if !ok {
		return op + ": " + beingName + " has no desk"
	}

	relationship := origin

	switch op {
	case "accept":
		if err := desk.AcceptTask(relationship, taskName, origin); err != nil {
			return "accept: " + err.Error()
		}
		debug.Log("[thread]: accepted", taskName, "on", beingName, "by", origin)
		return "accepted: " + taskName + " [" + beingName + "]"
	case "reject":
		if err := desk.RejectTask(relationship, taskName); err != nil {
			return "reject: " + err.Error()
		}
		debug.Log("[thread]: rejected", taskName, "on", beingName, "by", origin)
		return "rejected: " + taskName + " [" + beingName + "] — reopened"
	}
	return ""
}

func (t *NewThread) Grow(impulse string) string {
	name, err := Extract(impulse, "~name", "being")
	if err != nil {
		return "being: missing ~name"
	}
	if _, exists := t.Beings[name]; exists {
		return "being: " + name + " already exists"
	}
	beingType, err := Extract(impulse, "~type", "being")
	if err != nil {
		return "being: missing ~type"
	}
	devicesRaw, err := Extract(impulse, "~devices", "being")
	if err != nil {
		return "being: missing ~devices"
	}

	ctx := &Relation{ID: name, Impulse: impulse}
	ctx.Realities = make(map[string]Reality)

	for _, devName := range strings.Split(devicesRaw, ",") {
		devName = strings.TrimSpace(devName)
		if dev, ok := t.Devices[devName]; ok {
			ctx.Realities[devName] = dev
			if mac, ok := dev.(*MacOS); ok {
				for compName, comp := range mac.Components {
					ctx.Realities[compName] = comp
				}
			}
		}
	}

	var created Reality
	switch beingType {
	case "llm":
		created = (&Self{}).Create(ctx)
	case "user":
		created = (&User{}).Create(ctx)
		t.Access[name] = true
	case "cli":
		created = (&CLI{}).Create(ctx)
	case "agent":
		created = (&Agent{}).Create(ctx)
	default:
		return "being: unknown type " + beingType
	}

	t.Beings[name] = created

	being := extractBeing(created)

	for _, peerName := range being.Relationships {
		if peer, ok := t.Beings[peerName]; ok {
			peerBeing := extractBeing(peer)
			peerBeing.Relationships = append(peerBeing.Relationships, name)
			setBeing(peer, peerBeing)
		}
	}

	debug.Log("[thread]: created", name, "type:", beingType, "devices:", devicesRaw)
	return "created " + name
}

func extractBeing(r Reality) Being {
	switch c := r.(type) {
	case *Self:
		if b, ok := c.Realities["being"].(Being); ok {
			return b
		}
	case *User:
		if b, ok := c.Realities["being"].(Being); ok {
			return b
		}
	case *CLI:
		if b, ok := c.Realities["being"].(Being); ok {
			return b
		}
	case *Agent:
		if b, ok := c.Realities["being"].(Being); ok {
			return b
		}
	}
	return Being{}
}

func setBeing(r Reality, b Being) {
	switch c := r.(type) {
	case *Self:
		c.Realities["being"] = b
	case *User:
		c.Realities["being"] = b
	case *CLI:
		c.Realities["being"] = b
	case *Agent:
		c.Realities["being"] = b
	}
}
