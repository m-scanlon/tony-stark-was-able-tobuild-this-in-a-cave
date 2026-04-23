package thread

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
	"skyra-v04/src/primitives/operator"
)

var _ operator.IOperator = (*CloseThread)(nil)

type CloseThread struct {
	presentThread
	EntityMap map[string]entity.Entity
}

func (c *CloseThread) Relate(r entity.Relation) entity.Entity {
	id, _ := meaning.Extract(r.Impulse, "~with", "close-thread")
	if t, ok := c.EntityMap[id].(Thread); ok {
		t.Active = false
		c.EntityMap[id] = t
	}
	return c
}

func (c *CloseThread) ID() string   { return "close-thread" }
func (c *CloseThread) Name() string { return "close-thread" }
