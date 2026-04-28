// Being is a world of two: an inner entity and an outer entity. The inner
// entity deliberates first, the outer entity speaks. Both resolve through
// an LLM world of inference providers.
package world

import (
	"skyra-v05/src/entity"
	"skyra-v05/src/entity/being"
)

type ExchangePair struct {
	Thought string
	Output  string
}

type Being struct {
	World
	Pathos  being.Being
	Inner   entity.Entity
	Outer   entity.Entity
	LLM     *LLM
	Window  []ExchangePair
}

func NewBeing(pathos being.Being, llm *LLM) *Being {
	return &Being{
		World: World{
			Entities: make(map[string]entity.Entity),
		},
		Pathos: pathos,
		LLM:    llm,
		Window: make([]ExchangePair, 0, 10),
	}
}

func (b *Being) DerivePresent(r entity.Relation) string {
	return ""
}
