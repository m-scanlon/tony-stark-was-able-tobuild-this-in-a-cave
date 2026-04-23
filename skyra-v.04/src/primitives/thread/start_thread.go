package thread

import (
	"crypto/rand"
	"fmt"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
	"skyra-v04/src/primitives/operator"
)

var _ operator.IOperator = (*StartThread)(nil)

type StartThread struct {
	presentThread
	EntityMap map[string]entity.Entity
}

func newThreadID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (s *StartThread) Relate(r entity.Relation) entity.Entity {
	with, err := meaning.Extract(r.Impulse, "~with", "start-thread")
	if err != nil {
		fmt.Println("start-thread: missing ~with")
		return s
	}

	threadID := newThreadID()
	rel := r
	rel.ThreadID = threadID
	t, _ := Thread{}.Relate(rel).(Thread)
	s.EntityMap[threadID] = t

	// Kick off: send the initiator's opening message to the target via continue-thread.
	ct, ok := s.EntityMap["continue-thread"]
	if !ok {
		return t
	}
	say, _ := meaning.Extract(r.Impulse, "~say", "start-thread", "|")
	initRaw := fmt.Sprintf("skyra continue-thread ~with %s ~say %s | start", with, say)
	initRel, err := entity.Impress(r.Origin, threadID, initRaw)
	if err != nil {
		return t
	}
	ct.Relate(initRel)
	return t
}

func (s *StartThread) ID() string   { return "start-thread" }
func (s *StartThread) Name() string { return "start-thread" }
