package thread

import (
	"skyra-v04/src/primitives/entity"
)

type StartThread struct {
	presentThread
	EntityMap map[string]entity.Entity
}

func (s *StartThread) Relate(r entity.Relation) entity.Entity {
	thread, _ := Thread{}.Relate(r).(Thread)
	target, ok := s.EntityMap[r.ID]
	if !ok {
		return s
	}
	target.Relate(entity.Relation{
		ID:       thread.id,
		Origin:   r.Origin,
		ThreadID: thread.id,
		Impulse:  r.Impulse,
	})
	return thread
}

func (s *StartThread) ID() string   { return "start-thread" }
func (s *StartThread) Name() string { return "start-thread" }
