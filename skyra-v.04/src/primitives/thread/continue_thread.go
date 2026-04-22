package thread

import (
	"fmt"

	"skyra-v04/src/inference"
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

type inferrable interface {
	Name() string
	DerivePresent() string
	Receive(origin, entry string) logos.Logos
}

type ContinueThread struct {
	LogosMap map[string]logos.Logos
}

func (c *ContinueThread) Relate(r logos.Relation) logos.Logos {
	name, _ := meaning.Extract(r.Impulse, "~with", "continue-thread")
	target, ok := c.LogosMap[name]
	if !ok {
		return c
	}
	b, ok := target.(inferrable)
	if !ok {
		return c
	}

	// target's exchange with origin — arrival and response (same slot)
	updated := b.Receive(r.Origin, r.Impulse)
	c.LogosMap[name] = updated

	// origin's exchange with target — inbound and response (same slot, if origin is a being)
	if origin, ok := c.LogosMap[r.Origin]; ok {
		if ob, ok := origin.(inferrable); ok {
			c.LogosMap[r.Origin] = ob.Receive(name, r.Impulse)
		}
	}

	present := updated.(inferrable).DerivePresent() + "\n" + r.Impulse
	response, err := inference.Call(present)
	if err != nil {
		fmt.Println("inference error:", err)
		return c
	}

	updated = updated.(inferrable).Receive(r.Origin, b.Name()+": "+response)
	c.LogosMap[name] = updated

	if origin, ok := c.LogosMap[r.Origin]; ok {
		if ob, ok := origin.(inferrable); ok {
			c.LogosMap[r.Origin] = ob.Receive(name, b.Name()+": "+response)
		}
	}

	fmt.Println(b.Name() + ": " + response)
	return c
}

func (c *ContinueThread) ID() string   { return "continue-thread" }
func (c *ContinueThread) Name() string { return "continue-thread" }
