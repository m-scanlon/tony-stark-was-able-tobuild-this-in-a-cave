package thread

import (
	"fmt"

	"skyra-v04/src/inference"
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

type inferrable interface {
	entity.Entity
	Name() string
	Receive(origin, entry string) entity.Entity
}

type ContinueThread struct {
	presentThread
	EntityMap map[string]entity.Entity
}

func traceRelation(stage string, r entity.Relation) {
	fmt.Printf("trace: %s origin=%s target=%s thread=%s impulse=%q\n", stage, r.Origin, r.ID, r.ThreadID, r.Impulse)
}

func (c *ContinueThread) Relate(r entity.Relation) entity.Entity {
	traceRelation("ingress", r)
	return c.relate(r)
}

func (c *ContinueThread) relate(r entity.Relation) entity.Entity {
	name, _ := meaning.Extract(r.Impulse, "~with", "continue-thread")
	message := c.DerivePresent(r)
	target, ok := c.EntityMap[name]
	if !ok {
		fmt.Println("debug: target not found:", name)
		return c
	}
	b, ok := target.(inferrable)
	if !ok {
		fmt.Println("debug: target not inferrable:", name)
		return c
	}

	// target's exchange with origin — arrival and response (same slot)
	updated := b.Receive(r.Origin, message)
	c.EntityMap[name] = updated

	// origin's exchange with target — inbound and response (same slot, if origin is a being)
	if origin, ok := c.EntityMap[r.Origin]; ok {
		if ob, ok := origin.(inferrable); ok {
			c.EntityMap[r.Origin] = ob.Receive(name, message)
		}
	}

	senderContext := ""
	if origin, ok := c.EntityMap[r.Origin]; ok {
		if ob, ok := origin.(inferrable); ok {
			senderContext = "\nsender:\n" + ob.DerivePresent(r)
		}
	}

	beingsContext := "\nbeings you can talk to:\n"
	for id, node := range c.EntityMap {
		if id == name || id == r.Origin {
			continue
		}
		if nb, ok := node.(inferrable); ok {
			beingsContext += "  " + nb.Name() + "\n"
		}
	}

	present := updated.(inferrable).DerivePresent(r) + senderContext + beingsContext + "\nmessage from " + r.Origin + ": " + message
	response, err := inference.Call(present)
	if err != nil {
		fmt.Println("inference error:", err)
		return c
	}
	fmt.Println("debug: inference response received, len:", len(response))

	updated = updated.(inferrable).Receive(r.Origin, b.Name()+": "+response)
	c.EntityMap[name] = updated

	if origin, ok := c.EntityMap[r.Origin]; ok {
		if ob, ok := origin.(inferrable); ok {
			c.EntityMap[r.Origin] = ob.Receive(name, b.Name()+": "+response)
		}
	}

	if next, err := entity.Impress(name, r.ThreadID, response); err == nil {
		traceRelation("dispatch", next)
		if next.ID == r.Origin {
			fmt.Println(name + ": " + next.Impulse)
			return c
		}
		nextNode, ok := c.EntityMap[next.ID]
		if !ok {
			fmt.Println("debug: emitted target not found:", next.ID)
			return c
		}
		if next.ID == c.Name() {
			return c.relate(next)
		}
		nextNode.Relate(next)
		return c
	}

	fmt.Println(b.Name() + ": " + response)
	return c
}

func (c *ContinueThread) ID() string   { return "continue-thread" }
func (c *ContinueThread) Name() string { return "continue-thread" }
