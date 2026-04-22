package relationship

import (
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/thread"
)

type Relationship struct {
	peer    logos.Logos
	threads map[string]thread.Thread
}

func (r Relationship) Relate(rel logos.Relation) logos.Logos {
	t, _ := thread.Thread{}.Relate(rel).(thread.Thread)
	r.threads[t.ID()] = t
	return r
}

func (r Relationship) ID() string                    { return r.peer.ID() }
func (r Relationship) Name() string                  { return "relationship" }
func (r Relationship) Peer() logos.Logos             { return r.peer }
func (r Relationship) Threads() map[string]thread.Thread { return r.threads }
