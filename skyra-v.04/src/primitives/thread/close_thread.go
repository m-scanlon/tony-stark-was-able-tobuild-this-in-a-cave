package thread

import (
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

type CloseThread struct {
	LogosMap map[string]logos.Logos
}

func (c *CloseThread) Relate(r logos.Relation) logos.Logos {
	name, _ := meaning.Extract(r.Impulse, "~with", "close-thread")
	if t, ok := c.LogosMap[name]; ok {
		if thread, ok := t.(Thread); ok {
			thread.Active = false
			c.LogosMap[name] = thread
		}
	}
	return c
}

func (c *CloseThread) ID() string   { return "close-thread" }
func (c *CloseThread) Name() string { return "close-thread" }
