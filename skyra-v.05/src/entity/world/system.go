// System is the top-level world. It contains being worlds, routes messages
// between them, and manages threads and exchanges.
package world

import (
	"skyra-v05/src/entity"
)

type System struct {
	World
	threads map[string]*Thread
}

func NewSystem() *System {
	return &System{
		World: World{
			Entities: make(map[string]entity.Entity),
		},
		threads: make(map[string]*Thread),
	}
}

func (s *System) DerivePresent(r entity.Relation) string {
	return ""
}
