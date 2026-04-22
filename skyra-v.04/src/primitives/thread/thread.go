package thread

import (
	"skyra-v04/src/primitives/exchange"
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

type Thread struct {
	id        string
	About     string
	Because   string
	Active    bool
	exchanges map[string]exchange.Exchange
}

func (t Thread) Relate(r logos.Relation) logos.Logos {
	about, _ := meaning.Extract(r.Impulse, "~about", "thread")
	because, _ := meaning.Extract(r.Impulse, "~because", "thread")
	return Thread{
		id:        r.ThreadID,
		About:     about,
		Because:   because,
		Active:    true,
		exchanges: make(map[string]exchange.Exchange),
	}
}

func (t Thread) ID() string   { return t.id }
func (t Thread) Name() string { return "thread" }

func (t Thread) AllEntries() []string {
	var all []string
	for _, ex := range t.exchanges {
		all = append(all, ex.Entries()...)
	}
	return all
}
