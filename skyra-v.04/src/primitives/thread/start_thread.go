package thread

import (
	"skyra-v04/src/primitives/logos"
)

type StartThread struct {
	LogosMap map[string]logos.Logos
}

func (s *StartThread) Relate(r logos.Relation) logos.Logos {
	thread, _ := Thread{}.Relate(r).(Thread)
	target, ok := s.LogosMap[r.ID]
	if !ok {
		return s
	}
	target.Relate(logos.Relation{
		ID:       thread.id,
		Origin:   r.Origin,
		ThreadID: thread.id,
		Impulse:  r.Impulse,
	})
	return thread
}

func (s *StartThread) ID() string   { return "start-thread" }
func (s *StartThread) Name() string { return "start-thread" }
