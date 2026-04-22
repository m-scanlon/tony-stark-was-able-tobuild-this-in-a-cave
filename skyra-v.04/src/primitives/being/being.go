package being

import (
	"strings"

	"skyra-v04/src/primitives/impression"
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/nature"
	"skyra-v04/src/primitives/relationship"
)

type IBeing interface {
	logos.Logos
	DerivePresent() string
}

type Being struct {
	id            string
	name          string
	Impression    string
	nature        nature.Nature
	relationships map[string]relationship.Relationship
	exchanges     map[string][]string // keyed by peer name
}

func (b Being) ID() string   { return b.id }
func (b Being) Name() string { return b.name }

func (b Being) Relate(r logos.Relation) logos.Logos {
	n, _ := nature.Nature{}.Relate(r).(nature.Nature)
	imp, _ := impression.Impression{}.Relate(r).(impression.Impression)
	return Being{
		id:            strings.TrimSpace(r.ID),
		name:          strings.TrimSpace(r.ID),
		Impression:    imp.Value,
		nature:        n,
		relationships: make(map[string]relationship.Relationship),
		exchanges:     make(map[string][]string),
	}
}

func (b Being) Receive(origin, entry string) logos.Logos {
	newExchanges := make(map[string][]string, len(b.exchanges))
	for k, v := range b.exchanges {
		newExchanges[k] = v
	}
	newExchanges[origin] = append(newExchanges[origin], entry)
	b.exchanges = newExchanges
	return b
}

func (b Being) DerivePresent() string {
	var sb strings.Builder

	sb.WriteString("being: " + b.name + "\n")
	if b.nature.Identity.Value != "" {
		sb.WriteString("identity: " + b.nature.Identity.Value + "\n")
	}
	if b.nature.Purpose.Value != "" {
		sb.WriteString("purpose: " + b.nature.Purpose.Value + "\n")
	}
	if b.Impression != "" {
		sb.WriteString("impression: " + b.Impression + "\n")
	}

	if len(b.exchanges) > 0 {
		sb.WriteString("\nexchanges:\n")
		for peer, entries := range b.exchanges {
			sb.WriteString("  with: " + peer + "\n")
			for _, entry := range entries {
				sb.WriteString("    - " + entry + "\n")
			}
		}
	}

	sb.WriteString("\noperators:\n")
	sb.WriteString("  start-thread    ~about <topic> ~because <reason>\n")
	sb.WriteString("  continue-thread ~with <thread-id> <message>\n")
	sb.WriteString("  close-thread    ~with <thread-id>\n")

	return sb.String()
}
