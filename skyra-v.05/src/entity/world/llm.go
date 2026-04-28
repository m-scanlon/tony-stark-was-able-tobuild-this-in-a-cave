// LLM is a world of inference providers. Its DerivePresent selects a provider
// and routes the present to it. Each provider is an invariant.
package world

import (
	"skyra-v05/src/entity"
)

type LLM struct {
	World
}

func NewLLM() *LLM {
	return &LLM{
		World: World{
			Entities: make(map[string]entity.Entity),
		},
	}
}

func (l *LLM) DerivePresent(r entity.Relation) string {
	return ""
}
