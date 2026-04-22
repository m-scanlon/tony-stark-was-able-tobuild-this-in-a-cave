package thread

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

type CloseThread struct {
	presentThread
	EntityMap map[string]entity.Entity
}

func (c *CloseThread) Relate(r entity.Relation) entity.Entity {
	name, _ := meaning.Extract(r.Impulse, "~with", "close-thread")
	if t, ok := c.EntityMap[name]; ok {
		if thread, ok := t.(Thread); ok {
			thread.Active = false
			c.EntityMap[name] = thread
		}
	}
	return c
}

func (c *CloseThread) ID() string   { return "close-thread" }
func (c *CloseThread) Name() string { return "close-thread" }
