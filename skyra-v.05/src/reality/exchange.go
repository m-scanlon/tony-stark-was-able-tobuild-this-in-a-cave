package reality

import (
	"fmt"
	"skyra-v05/src/debug"
	"strconv"
	"strings"
	"time"
)

type Exchange struct {
	id        string
	Exchanges map[string]*Conversation
	Levels    *Levels
}

type Conversation struct {
	Parties [2]string
	Active  bool
	Entries []Entry
	Sender  string
	Message string
	Context map[string]string
}

type Entry struct {
	From    string
	Content string
	Time    time.Time
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
	if r.Collecting {
		for key, conv := range e.Exchanges {
			snap := ExchangeSnapshot{
				Key:     key,
				Parties: conv.Parties,
				Active:  conv.Active,
				Entries: []EntrySnapshot{},
			}
			for i, entry := range conv.Entries {
				snap.Entries = append(snap.Entries, EntrySnapshot{
					Index: i, From: entry.From,
					Content: entry.Content, Ts: entry.Time.UnixMilli(),
				})
			}
			if len(conv.Context) > 0 {
				snap.Context = conv.Context
			}
			r.Export("exchange:"+key, snap)
		}
		r.Export("node:exchange", RealityNode{ID: "exchange", Type: "Exchange", Children: []RealityNode{}})
		return ""
	}

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

	refValue, refErr := ExtractTag(r.Impulse, "ref")
	if refErr == nil {
		debug.Log("[exchange]: <ref> found →", refValue)
		r.Impulse = StripTag(r.Impulse, "ref")

		srcPeer, start, end, parseErr := parseRef(refValue)
		if parseErr == nil {
			srcKey := exchangeKey(r.Origin, srcPeer)
			if srcConv, ok := e.Exchanges[srcKey]; ok {
				resolved := srcConv.SliceEntries(start, end)
				debug.Log("[exchange]: resolved ~ref", srcKey, start, "-", end, "→", len(resolved), "entries")

				srcConv.Entries = append(srcConv.Entries, Entry{
					From:    r.Origin,
					Content: fmt.Sprintf("[left to talk to %s]", r.ID),
					Time:    time.Now(),
				})
				debug.Log("[exchange]: appended departure to", srcKey)

				if conv.Context == nil {
					conv.Context = make(map[string]string)
				}
				conv.Context[r.Origin] = renderRef(srcPeer, start, end, resolved)
				debug.Log("[exchange]: stored context for", r.Origin, "on", key)
			} else {
				debug.Log("[exchange]: ~ref source conversation not found:", srcKey)
			}
		} else {
			debug.Log("[exchange]: ~ref parse error:", parseErr)
		}
	}

	/*
		Ref enforcement is disabled for the current experiment.
		Messages may cross from one exchange to another without carrying a
		<ref> block. Explicit refs are still parsed above when present.

		isProcess := false
		if target, ok := r.Realities[r.ID]; ok {
			if _, ok := target.(*Process); ok {
				isProcess = true
			}
		}

		if !hasRef && !ok && !isProcess {
			for existingKey, existingConv := range e.Exchanges {
				if existingConv.Active && existingKey != key {
					for _, party := range existingConv.Parties {
						if party == r.Origin {
							otherPeer := existingConv.Parties[0]
							if otherPeer == r.Origin {
								otherPeer = existingConv.Parties[1]
							}
							debug.Log("[exchange]: crossing without ~ref — blocking", r.Origin, "→", r.ID)
							delete(e.Exchanges, key)
							r.Realities["error"] = &Error{
								Message: fmt.Sprintf("you are leaving your exchange with %s to talk to %s without carrying context. use <ref>%s:START-END</ref> inside your tag to bring context from that conversation.", otherPeer, r.ID, otherPeer),
							}
							return ""
						}
					}
				}
			}
		}
	*/

	conv.Sender = r.Origin
	conv.Message = r.Impulse
	debug.Log("[exchange]: recording entry from", r.Origin, "→", r.Impulse[:min(len(r.Impulse), 40)])
	conv.Entries = append(conv.Entries, Entry{
		From:    r.Origin,
		Content: r.Impulse,
		Time:    time.Now(),
	})

	if e.Levels != nil {
		e.Levels.Award(r.Origin, r.ID)
	}

	const compactThreshold = 20
	const keepRecent = 10
	if len(conv.Entries) > compactThreshold {
		if mem := findMemory(r); mem != nil {
			older := conv.Entries[:len(conv.Entries)-keepRecent]
			relationship := r.Origin
			if relationship == mem.Owner {
				relationship = r.ID
			}
			debug.Log("[exchange]: compacting", len(older), "entries to memory for", relationship)
			mem.Compress(older, relationship)
			conv.Entries = conv.Entries[len(conv.Entries)-keepRecent:]
		}
	}
	r.Attach("exchange", conv.Parse)
	r.Attach("conversation", func() string {
		return conv.ParseRecent(10)
	})
	target := r.ID
	if ctx := conv.ContextFor(target); ctx != "" {
		r.Attach("ref-context", func() string {
			return ctx
		})
	}

	being, ok := r.Realities[r.ID]
	if !ok {
		debug.Log("[exchange]: being not found:", r.ID, "| available:", len(r.Realities))
		for name := range r.Realities {
			debug.Log("[exchange]:   -", name)
		}
		return ""
	}

	if proc, ok := being.(*Process); ok {
		debug.Log("[exchange]: direct process call:", r.ID)
		result := proc.Realize(r)
		conv.Entries = append(conv.Entries, Entry{
			From:    r.ID,
			Content: result,
			Time:    time.Now(),
		})
		r.Origin = r.ID
		r.ID = conv.Parties[0]
		if r.ID == r.Origin {
			r.ID = conv.Parties[1]
		}
		r.Impulse = result
		return result
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

func (c *Conversation) ID() string { return exchangeKey(c.Parties[0], c.Parties[1]) }

func (c *Conversation) Create(r *Relation) Reality { return c }

func (c *Conversation) Realize(r *Relation) string { return "" }

func (c *Conversation) ParseRecent(n int) string {
	entries := c.Entries
	if len(entries) > n {
		entries = entries[len(entries)-n:]
	}
	var sb strings.Builder
	sb.WriteString("recent conversation:\n")
	for _, entry := range entries {
		sb.WriteString("  " + entry.From + ": " + entry.Content + "\n")
	}
	return sb.String()
}

func parseRef(value string) (string, int, int, error) {
	peer, rang, ok := strings.Cut(value, ":")
	if !ok {
		return "", 0, 0, fmt.Errorf("ref: missing colon in %q", value)
	}
	startStr, endStr, hasRange := strings.Cut(rang, "-")
	start, err := strconv.Atoi(strings.TrimSpace(startStr))
	if err != nil {
		return "", 0, 0, fmt.Errorf("ref: bad start %q", startStr)
	}
	if !hasRange {
		return strings.TrimSpace(peer), start, start, nil
	}
	end, err := strconv.Atoi(strings.TrimSpace(endStr))
	if err != nil {
		return "", 0, 0, fmt.Errorf("ref: bad end %q", endStr)
	}
	return strings.TrimSpace(peer), start, end, nil
}

func (c *Conversation) SliceEntries(start, end int) []Entry {
	if start < 0 {
		start = 0
	}
	if end >= len(c.Entries) {
		end = len(c.Entries) - 1
	}
	if start > end || start >= len(c.Entries) {
		return nil
	}
	return c.Entries[start : end+1]
}

func renderRef(peer string, start, end int, entries []Entry) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("context brought from %s (entries %d-%d):\n", peer, start, end))
	for i, entry := range entries {
		sb.WriteString(fmt.Sprintf("  [%d] %s: %s\n", start+i, entry.From, entry.Content))
	}
	return sb.String()
}

func (c *Conversation) ContextFor(being string) string {
	if c.Context == nil {
		return ""
	}
	return c.Context[being]
}
