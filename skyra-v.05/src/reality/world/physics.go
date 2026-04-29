package world

import (
	"fmt"
	"sort"
	"strconv"

	"skyra-v05/src/reality"
	"skyra-v05/src/reality/being"
)

var physicsRegistry = map[string]func() reality.Reality{
	"thread":    func() reality.Reality { return NewThread() },
	"economics": func() reality.Reality { return NewEconomics() },
}

func RegisterPhysics(name string, constructor func() reality.Reality) {
	physicsRegistry[name] = constructor
}

type physicsEntry struct {
	order int
	reality reality.Reality
}

type Physics struct {
	id      string
	entries []physicsEntry
	lookup  map[string]reality.Reality
}

func NewPhysics() *Physics {
	return &Physics{
		id:     "physics",
		lookup: make(map[string]reality.Reality),
	}
}

func (p *Physics) ID() string { return p.id }

func (p *Physics) Create(r reality.Relation) reality.Reality {
	name, _ := being.Extract(r.Impulse, "~name", "physics")
	orderStr, _ := being.Extract(r.Impulse, "~order", "physics")
	order, _ := strconv.Atoi(orderStr)

	constructor, ok := physicsRegistry[name]
	if !ok {
		fmt.Println("physics: unknown type:", name)
		return p
	}

	entry := physicsEntry{order: order, reality: constructor()}
	p.entries = append(p.entries, entry)
	p.lookup[name] = entry.reality

	sort.Slice(p.entries, func(i, j int) bool {
		return p.entries[i].order < p.entries[j].order
	})

	return p
}

func (p *Physics) Realize(r reality.Relation) string {
	result := ""
	for _, entry := range p.entries {
		result += entry.reality.Realize(r)
	}
	return result
}

func (p *Physics) Get(name string) reality.Reality {
	return p.lookup[name]
}
