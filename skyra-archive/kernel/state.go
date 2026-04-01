package kernel

import (
	"encoding/json"
	"sync"
	"time"
)

type InteractionSnapshot struct {
	Messages []InteractionMessage `json:"messages"`
}

const staticPerceptionID = "perception_v1"

type ThoughtSnapshot struct {
	Perception  *Perception         `json:"perception,omitempty"`
	ActiveStep  *ThoughtStepState   `json:"active_step,omitempty"`
	ActiveChain *ChainState         `json:"active_chain,omitempty"`
	Suspended   []ChainState        `json:"suspended_chains,omitempty"`
	Steps       []ThoughtStepRecord `json:"steps"`
	QueueDepth  int                 `json:"queue_depth"`
}

type runtimeState struct {
	mu                sync.RWMutex
	stimuli           []Stimulus
	history           []HistoryEntry
	interaction       []InteractionMessage
	steps             []ThoughtStepRecord
	activeStep        *ThoughtStepState
	activeChain       *ChainState
	suspended         []*ChainState
	perception        *Perception
	lastUnderstanding *Understanding
	nextStepIndex     int
	historyLimit      int
}

func newRuntimeState(historyLimit int) *runtimeState {
	if historyLimit <= 0 {
		historyLimit = 12
	}
	now := time.Now().UTC()
	return &runtimeState{
		historyLimit: historyLimit,
		perception: &Perception{
			ID:        staticPerceptionID,
			UpdatedAt: now,
		},
	}
}

func (s *runtimeState) appendStimulus(source, stimulusType, rawPayload string, now time.Time) Stimulus {
	s.mu.Lock()
	defer s.mu.Unlock()

	stimulus := Stimulus{
		ID:                newID(),
		Type:              stimulusType,
		Source:            source,
		RawPayload:        rawPayload,
		NormalizedSummary: rawPayload,
		Timestamp:         now,
	}
	s.stimuli = append(s.stimuli, stimulus)
	return stimulus
}

func (s *runtimeState) appendHistory(role, source, content string, now time.Time) HistoryEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := HistoryEntry{
		ID:        newID(),
		Source:    source,
		Role:      role,
		Content:   content,
		Timestamp: now,
	}
	s.history = append(s.history, entry)
	if len(s.history) > s.historyLimit {
		s.history = append([]HistoryEntry(nil), s.history[len(s.history)-s.historyLimit:]...)
	}
	return entry
}

func (s *runtimeState) appendInteraction(role, content, stimulusID string, now time.Time) InteractionMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := InteractionMessage{
		ID:         newID(),
		Role:       role,
		Content:    content,
		StimulusID: stimulusID,
		CreatedAt:  now,
	}
	s.interaction = append(s.interaction, msg)
	return msg
}

func (s *runtimeState) startChain(stimulus Stimulus, now time.Time) ChainState {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeChain != nil {
		s.activeChain.Status = "superseded"
		s.activeChain.UpdatedAt = now
	}

	history := append([]HistoryEntry(nil), s.history...)
	perception := clonePerception(s.perception)
	if perception == nil {
		perception = &Perception{ID: staticPerceptionID}
	}
	perception.ID = staticPerceptionID
	perception.History = history
	perception.Stimulus = stimulus
	perception.UpdatedAt = now

	chain := ChainState{
		ID:               newID(),
		Version:          1,
		Status:           "running",
		RootStimulusID:   stimulus.ID,
		StimulusSource:   stimulus.Source,
		ReviewingUser:    isHumanSource(stimulus.Source) && len(s.suspended) > 0,
		CurrentPrimitive: string(EventKindExperience),
		Perception:       *perception,
		UpdatedAt:        now,
	}

	s.activeChain = cloneChainState(&chain)
	s.perception = clonePerception(perception)
	return chain
}

func (s *runtimeState) activeChainForEvent(event *Event) (*ChainState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if event == nil || s.activeChain == nil {
		return nil, false
	}
	if s.activeChain.ID != event.ChainID || s.activeChain.Version != event.ChainVersion {
		return nil, false
	}
	if s.activeChain.Status == "completed" || s.activeChain.Status == "superseded" || s.activeChain.Status == "failed" {
		return nil, false
	}
	return cloneChainState(s.activeChain), true
}

func (s *runtimeState) currentUnderstanding() *Understanding {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.lastUnderstanding == nil {
		return nil
	}
	copyValue := *s.lastUnderstanding
	return &copyValue
}

func (s *runtimeState) hasSuspendedChains() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.suspended) > 0
}

func (s *runtimeState) activeChainSnapshot() *ChainState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneChainState(s.activeChain)
}

func (s *runtimeState) activeStepSnapshot(stepID string) (*ThoughtStepState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.activeStep == nil {
		return nil, false
	}
	if stepID != "" && s.activeStep.StepID != stepID {
		return nil, false
	}
	return cloneThoughtStepState(s.activeStep), true
}

func (s *runtimeState) beginStep(chainID string, version int, frame string, now time.Time) (ThoughtStepState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeChain == nil || s.activeChain.ID != chainID || s.activeChain.Version != version {
		return ThoughtStepState{}, false
	}

	s.nextStepIndex++
	s.activeChain.CurrentPrimitive = frame
	s.activeChain.Status = "running"
	s.activeChain.UpdatedAt = now

	step := ThoughtStepState{
		ChainID:   chainID,
		StepID:    newID(),
		StepIndex: s.nextStepIndex,
		Frame:     frame,
		Status:    "streaming",
		StartedAt: now,
	}
	s.activeStep = cloneThoughtStepState(&step)
	return step, true
}

func (s *runtimeState) appendStepDelta(stepID, delta string) *ThoughtStepState {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeStep == nil || s.activeStep.StepID != stepID {
		return nil
	}
	s.activeStep.RawOutput += delta
	copyStep := *s.activeStep
	return &copyStep
}

func (s *runtimeState) completeStep(step ThoughtStepState, status, primitiveChoice string, perception Perception, interactionChannel, interactionMessage, currentUnderstanding string, now time.Time) ThoughtStepRecord {
	s.mu.Lock()
	defer s.mu.Unlock()

	step.Status = status
	step.PrimitiveChoice = primitiveChoice
	step.CompletedAt = &now

	record := ThoughtStepRecord{
		ChainID:              step.ChainID,
		StepID:               step.StepID,
		StepIndex:            step.StepIndex,
		Frame:                step.Frame,
		Status:               status,
		PrimitiveChoice:      primitiveChoice,
		RawOutput:            step.RawOutput,
		PerceptionSnapshot:   perception,
		StartedAt:            step.StartedAt,
		CompletedAt:          now,
		InteractionChannel:   interactionChannel,
		InteractionMessage:   interactionMessage,
		CurrentUnderstanding: currentUnderstanding,
	}
	s.steps = append(s.steps, record)
	s.activeStep = nil
	return cloneThoughtStepRecord(record)
}

func (s *runtimeState) updateActiveChain(chainID string, version int, fn func(*ChainState)) (*ChainState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeChain == nil || s.activeChain.ID != chainID || s.activeChain.Version != version {
		return nil, false
	}
	fn(s.activeChain)
	s.activeChain.UpdatedAt = time.Now().UTC()
	s.activeChain.Perception.ID = staticPerceptionID
	s.perception = clonePerception(&s.activeChain.Perception)
	if s.activeChain.Perception.Understanding != nil {
		understanding := *s.activeChain.Perception.Understanding
		s.lastUnderstanding = &understanding
	}
	return cloneChainState(s.activeChain), true
}

func (s *runtimeState) completeActiveChain(chainID string, version int, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeChain == nil || s.activeChain.ID != chainID || s.activeChain.Version != version {
		return
	}
	s.activeChain.Status = "completed"
	s.activeChain.UpdatedAt = now
	s.activeChain = nil
}

func (s *runtimeState) suspendActiveChain(chainID string, version int, now time.Time) (*ChainState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeChain == nil || s.activeChain.ID != chainID || s.activeChain.Version != version {
		return nil, false
	}
	if s.activeChain.Status == "completed" || s.activeChain.Status == "failed" || s.activeChain.Status == "suspended" {
		return nil, false
	}

	s.activeChain.Status = "suspended"
	s.activeChain.UpdatedAt = now
	suspended := cloneChainState(s.activeChain)
	s.suspended = append(s.suspended, suspended)
	s.activeChain = nil
	return cloneChainState(suspended), true
}

func (s *runtimeState) suspendCurrentActiveChain(now time.Time) (*ChainState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeChain == nil {
		return nil, false
	}
	if s.activeChain.Status == "completed" || s.activeChain.Status == "failed" || s.activeChain.Status == "suspended" {
		return nil, false
	}

	s.activeChain.Status = "suspended"
	s.activeChain.UpdatedAt = now
	suspended := cloneChainState(s.activeChain)
	s.suspended = append(s.suspended, suspended)
	s.activeChain = nil
	return cloneChainState(suspended), true
}

func (s *runtimeState) resumeSuspendedChain(now time.Time) (*ChainState, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.suspended) == 0 {
		return nil, false
	}

	suspended := cloneChainState(s.suspended[len(s.suspended)-1])
	s.suspended = s.suspended[:len(s.suspended)-1]

	perception := clonePerception(s.perception)
	if perception == nil {
		perception = &Perception{ID: staticPerceptionID}
	}
	perception.ID = staticPerceptionID
	perception.History = append([]HistoryEntry(nil), s.history...)
	perception.Stimulus = suspended.Perception.Stimulus
	perception.UpdatedAt = now

	suspended.Version++
	suspended.Status = "running"
	suspended.CurrentPrimitive = string(EventKindReference)
	suspended.ReferenceOutput = ""
	suspended.InferOutput = ""
	suspended.InferLoopCount = 0
	suspended.Perception = *perception
	suspended.UpdatedAt = now

	s.activeChain = cloneChainState(suspended)
	s.perception = clonePerception(perception)
	if perception.Understanding != nil {
		understanding := *perception.Understanding
		s.lastUnderstanding = &understanding
	}
	return cloneChainState(s.activeChain), true
}

func (s *runtimeState) failActiveChain(chainID string, version int, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeChain == nil || s.activeChain.ID != chainID || s.activeChain.Version != version {
		return
	}
	s.activeChain.Status = "failed"
	s.activeChain.UpdatedAt = now
}

func (s *runtimeState) interactionSnapshot() InteractionSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]InteractionMessage, len(s.interaction))
	copy(out, s.interaction)
	return InteractionSnapshot{Messages: out}
}

func (s *runtimeState) thoughtSnapshot(queueDepth int) ThoughtSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	steps := make([]ThoughtStepRecord, len(s.steps))
	copy(steps, s.steps)

	suspended := make([]ChainState, 0, len(s.suspended))
	for _, chain := range s.suspended {
		if chain == nil {
			continue
		}
		suspended = append(suspended, *cloneChainState(chain))
	}

	return ThoughtSnapshot{
		Perception:  clonePerception(s.perception),
		ActiveStep:  cloneThoughtStepState(s.activeStep),
		ActiveChain: cloneChainState(s.activeChain),
		Suspended:   suspended,
		Steps:       steps,
		QueueDepth:  queueDepth,
	}
}

func (s *runtimeState) formatHistoryJSON() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return mustJSON(s.history)
}

func clonePerception(value *Perception) *Perception {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.History = append([]HistoryEntry(nil), value.History...)
	if value.Understanding != nil {
		understanding := *value.Understanding
		copyValue.Understanding = &understanding
	}
	return &copyValue
}

func cloneChainState(value *ChainState) *ChainState {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.Perception = *clonePerception(&value.Perception)
	return &copyValue
}

func cloneThoughtStepState(value *ThoughtStepState) *ThoughtStepState {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}

func cloneThoughtStepRecord(value ThoughtStepRecord) ThoughtStepRecord {
	value.PerceptionSnapshot = *clonePerception(&value.PerceptionSnapshot)
	return value
}

func mustJSON(value any) string {
	buf, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(buf)
}
