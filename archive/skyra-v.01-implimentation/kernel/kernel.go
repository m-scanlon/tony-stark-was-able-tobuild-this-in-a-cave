package kernel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	defaultHistoryLimit   = 12
	maxReferenceInferHops = 4
	maxChainStepCount     = 24
	maxResolveCycles      = 6
	primitivePriority     = 50
)

var ErrStepInterrupted = errors.New("kernel: step interrupted")

type Config struct {
	HistoryLimit int
	PromptRoot   string
	GatewayWSURL string
	GatewayToken string
}

type Kernel struct {
	heap          Heap
	prompts       *PromptRegistry
	gateway       *GatewayClient
	state         *runtimeState
	thoughtBroker *Broker
	chatBroker    *Broker
	stepMu        sync.Mutex
	stepCancel    context.CancelFunc
}

func New(cfg Config, heap Heap) (*Kernel, error) {
	if heap == nil {
		return nil, fmt.Errorf("kernel: heap is not configured")
	}
	if strings.TrimSpace(cfg.PromptRoot) == "" {
		return nil, fmt.Errorf("kernel: prompt root is required")
	}
	if strings.TrimSpace(cfg.GatewayWSURL) == "" {
		return nil, fmt.Errorf("kernel: gateway websocket url is required")
	}

	prompts, err := LoadPromptRegistry(cfg.PromptRoot)
	if err != nil {
		return nil, err
	}

	historyLimit := cfg.HistoryLimit
	if historyLimit <= 0 {
		historyLimit = defaultHistoryLimit
	}

	return &Kernel{
		heap:          heap,
		prompts:       prompts,
		gateway:       NewGatewayClient(cfg.GatewayWSURL, cfg.GatewayToken),
		state:         newRuntimeState(historyLimit),
		thoughtBroker: NewBroker(),
		chatBroker:    NewBroker(),
	}, nil
}

func (k *Kernel) Run(ctx context.Context) {
	log.Println("kernel: execution loop started")
	for {
		event, err := k.heap.Pop(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("kernel: execution loop stopped")
				return
			}
			log.Printf("kernel: heap pop error: %v", err)
			continue
		}
		if event == nil {
			continue
		}
		log.Printf(
			"kernel: dequeued event=%s kind=%s priority=%d chain=%s version=%d stimulus=%s queue_depth=%d",
			event.ID,
			event.Kind,
			event.Priority,
			emptyDash(event.ChainID),
			event.ChainVersion,
			emptyDash(event.StimulusID),
			k.heap.Len(),
		)
		if err := k.processEvent(ctx, event); err != nil {
			if errors.Is(err, ErrStepInterrupted) {
				continue
			}
			log.Printf("kernel: event=%s kind=%s failed: %v", event.ID, event.Kind, err)
			if event.ChainID != "" {
				k.state.failActiveChain(event.ChainID, event.ChainVersion, time.Now().UTC())
			}
			now := time.Now().UTC()
			msg := k.state.appendInteraction("assistant", "Runtime error: "+err.Error(), event.StimulusID, now)
			k.chatBroker.Publish(UIEvent{Type: "interaction.message", Data: msg})
		}
	}
}

func (k *Kernel) Ready(ctx context.Context) error {
	return k.gateway.Ready(ctx)
}

func (k *Kernel) SubmitStimulus(ctx context.Context, source, stimulusType, payload string) (Stimulus, error) {
	now := time.Now().UTC()
	if source == "" {
		source = "unknown"
	}
	stimulus := k.state.appendStimulus(source, stimulusType, payload, now)

	role := "signal"
	priority := primitivePriority
	if isHumanSource(source) {
		role = "user"
		priority = 100
	}

	k.state.appendHistory(role, source, payload, now)
	if isHumanSource(source) {
		msg := k.state.appendInteraction("user", payload, stimulus.ID, now)
		k.chatBroker.Publish(UIEvent{Type: "interaction.message", Data: msg})
	}

	event := &Event{
		ID:         newID(),
		StimulusID: stimulus.ID,
		Source:     source,
		Kind:       EventKindStimulus,
		Priority:   priority,
		CreatedAt:  now,
	}

	if err := k.heap.Push(ctx, event); err != nil {
		return Stimulus{}, err
	}
	log.Printf(
		"kernel: accepted stimulus=%s source=%s type=%s user_source=%t priority=%d content=%q",
		stimulus.ID,
		source,
		stimulusType,
		isHumanSource(source),
		priority,
		summarizeForLog(payload, 140),
	)
	if isHumanSource(source) {
		log.Printf("kernel: interrupt requested by user stimulus=%s", stimulus.ID)
		k.interruptActiveStep()
	}
	return stimulus, nil
}

func (k *Kernel) InteractionSnapshot() InteractionSnapshot {
	return k.state.interactionSnapshot()
}

func (k *Kernel) ThoughtSnapshot() ThoughtSnapshot {
	return k.state.thoughtSnapshot(k.heap.Len())
}

func (k *Kernel) SubscribeInteraction() (int, <-chan UIEvent) {
	return k.chatBroker.Subscribe()
}

func (k *Kernel) UnsubscribeInteraction(id int) {
	k.chatBroker.Unsubscribe(id)
}

func (k *Kernel) SubscribeThoughts() (int, <-chan UIEvent) {
	return k.thoughtBroker.Subscribe()
}

func (k *Kernel) UnsubscribeThoughts(id int) {
	k.thoughtBroker.Unsubscribe(id)
}

func (k *Kernel) processEvent(ctx context.Context, event *Event) error {
	log.Printf(
		"kernel: processing event=%s kind=%s chain=%s version=%d stimulus=%s",
		event.ID,
		event.Kind,
		emptyDash(event.ChainID),
		event.ChainVersion,
		emptyDash(event.StimulusID),
	)
	switch event.Kind {
	case EventKindStimulus:
		return k.processStimulusEvent(ctx, event)
	case EventKindExperience:
		return k.processExperienceEvent(ctx, event)
	case EventKindUnderstand:
		return k.processUnderstandEvent(ctx, event)
	case EventKindReference:
		return k.processReferenceEvent(ctx, event)
	case EventKindInfer:
		return k.processInferEvent(ctx, event)
	case EventKindResolve:
		return k.processResolveEvent(ctx, event)
	case EventKindInteract:
		return k.processInteractEvent(ctx, event)
	default:
		return fmt.Errorf("unknown event kind %q", event.Kind)
	}
}

func (k *Kernel) processStimulusEvent(ctx context.Context, event *Event) error {
	stimulus, ok := k.lookupStimulus(event.StimulusID)
	if !ok {
		return fmt.Errorf("stimulus %q not found", event.StimulusID)
	}

	now := time.Now().UTC()
	if isHumanSource(stimulus.Source) {
		if suspended, ok := k.state.suspendCurrentActiveChain(now); ok {
			log.Printf(
				"kernel: suspended chain=%s version=%d current_primitive=%s due_to_stimulus=%s",
				suspended.ID,
				suspended.Version,
				emptyDash(suspended.CurrentPrimitive),
				stimulus.ID,
			)
			k.publishSuspendedChain(suspended)
		}
	}

	chain := k.state.startChain(stimulus, now)
	log.Printf(
		"kernel: started chain=%s version=%d reviewing_user=%t stimulus=%s source=%s content=%q",
		chain.ID,
		chain.Version,
		chain.ReviewingUser,
		stimulus.ID,
		stimulus.Source,
		summarizeForLog(stimulus.RawPayload, 140),
	)
	k.publishPerception(chain.Perception)
	k.publishChain(chain)
	return k.enqueuePrimitive(ctx, &chain, EventKindExperience)
}

func (k *Kernel) processExperienceEvent(ctx context.Context, event *Event) error {
	chain, ok := k.state.activeChainForEvent(event)
	if !ok {
		k.markStaleEvent(event, "active chain changed before experience")
		return nil
	}

	raw, _, err := k.runStep(ctx, *chain, "experience", k.mustRender("experience", map[string]string{
		"source":     chain.StimulusSource,
		"stimulus":   chain.Perception.Stimulus.RawPayload,
		"perception": mustJSON(chain.Perception),
	}))
	if err != nil {
		if errors.Is(err, ErrStepInterrupted) {
			return k.handleInterruptedChain(chain.ID, chain.Version)
		}
		return err
	}

	choice := parseChoice(raw, "understand", "interact")
	if choice == "" {
		choice = "understand"
	}
	updatedChain, ok := k.state.updateActiveChain(chain.ID, chain.Version, func(state *ChainState) {
		state.StepCount++
		state.CurrentPrimitive = choice
	})
	if !ok {
		k.markStaleEvent(event, "chain changed during experience completion")
		return nil
	}
	log.Printf(
		"kernel: experience chose=%s chain=%s version=%d stimulus=%s",
		choice,
		updatedChain.ID,
		updatedChain.Version,
		updatedChain.RootStimulusID,
	)
	k.completeActiveStep(choice, updatedChain.Perception, "", "", currentUnderstandingText(updatedChain.Perception.Understanding))

	if reachedChainLimit(updatedChain) {
		k.state.completeActiveChain(updatedChain.ID, updatedChain.Version, time.Now().UTC())
		k.publishChainState(nil)
		return nil
	}

	if updatedChain.ReviewingUser && updatedChain.ResolveCount > 0 && choice != "interact" && k.state.hasSuspendedChains() {
		log.Printf(
			"kernel: user review chain=%s completed without outward response, resuming suspended reasoning",
			updatedChain.ID,
		)
		k.state.completeActiveChain(updatedChain.ID, updatedChain.Version, time.Now().UTC())
		return k.resumeSuspendedChain(ctx)
	}

	switch choice {
	case "interact":
		return k.enqueuePrimitive(ctx, updatedChain, EventKindInteract)
	default:
		return k.enqueuePrimitive(ctx, updatedChain, EventKindUnderstand)
	}
}

func (k *Kernel) processUnderstandEvent(ctx context.Context, event *Event) error {
	chain, ok := k.state.activeChainForEvent(event)
	if !ok {
		k.markStaleEvent(event, "active chain changed before understand")
		return nil
	}

	raw, _, err := k.runStep(ctx, *chain, "understand", k.mustRender("understand", nil))
	if err != nil {
		if errors.Is(err, ErrStepInterrupted) {
			return k.handleInterruptedChain(chain.ID, chain.Version)
		}
		return err
	}
	_ = parseChoice(raw, "interpret")

	updatedChain, ok := k.state.updateActiveChain(chain.ID, chain.Version, func(state *ChainState) {
		state.StepCount++
		state.CurrentPrimitive = string(EventKindReference)
	})
	if !ok {
		k.markStaleEvent(event, "chain changed during understand completion")
		return nil
	}
	log.Printf(
		"kernel: understand advanced chain=%s version=%d next=%s",
		updatedChain.ID,
		updatedChain.Version,
		EventKindReference,
	)
	k.completeActiveStep("interpret", updatedChain.Perception, "", "", currentUnderstandingText(updatedChain.Perception.Understanding))

	if reachedChainLimit(updatedChain) {
		k.state.completeActiveChain(updatedChain.ID, updatedChain.Version, time.Now().UTC())
		k.publishChainState(nil)
		return nil
	}
	return k.enqueuePrimitive(ctx, updatedChain, EventKindReference)
}

func (k *Kernel) processReferenceEvent(ctx context.Context, event *Event) error {
	chain, ok := k.state.activeChainForEvent(event)
	if !ok {
		k.markStaleEvent(event, "active chain changed before reference")
		return nil
	}

	raw, _, err := k.runStep(ctx, *chain, "reference", k.mustRender("reference", map[string]string{
		"perception_frame": mustJSON(chain.Perception),
		"knowledge":        k.state.formatHistoryJSON(),
	}))
	if err != nil {
		if errors.Is(err, ErrStepInterrupted) {
			return k.handleInterruptedChain(chain.ID, chain.Version)
		}
		return err
	}
	referenceOutput := extractSection(raw, "current reference")
	if referenceOutput == "" {
		referenceOutput = strings.TrimSpace(raw)
	}

	updatedChain, ok := k.state.updateActiveChain(chain.ID, chain.Version, func(state *ChainState) {
		state.StepCount++
		state.CurrentPrimitive = string(EventKindInfer)
		state.ReferenceOutput = referenceOutput
	})
	if !ok {
		k.markStaleEvent(event, "chain changed during reference completion")
		return nil
	}
	log.Printf(
		"kernel: reference advanced chain=%s version=%d next=%s reference=%q",
		updatedChain.ID,
		updatedChain.Version,
		EventKindInfer,
		summarizeForLog(referenceOutput, 160),
	)
	k.completeActiveStep("infer", updatedChain.Perception, "", "", currentUnderstandingText(updatedChain.Perception.Understanding))

	if reachedChainLimit(updatedChain) {
		k.state.completeActiveChain(updatedChain.ID, updatedChain.Version, time.Now().UTC())
		k.publishChainState(nil)
		return nil
	}
	return k.enqueuePrimitive(ctx, updatedChain, EventKindInfer)
}

func (k *Kernel) processInferEvent(ctx context.Context, event *Event) error {
	chain, ok := k.state.activeChainForEvent(event)
	if !ok {
		k.markStaleEvent(event, "active chain changed before infer")
		return nil
	}

	raw, _, err := k.runStep(ctx, *chain, "infer", k.mustRender("infer", map[string]string{
		"perception_frame": mustJSON(chain.Perception),
		"reference_output": chain.ReferenceOutput,
	}))
	if err != nil {
		if errors.Is(err, ErrStepInterrupted) {
			return k.handleInterruptedChain(chain.ID, chain.Version)
		}
		return err
	}

	inferOutput := extractSection(raw, "current infer")
	if inferOutput == "" {
		inferOutput = strings.TrimSpace(raw)
	}
	nextChoice := parseChoice(extractSection(raw, "primitive choice"), "reference", "resolve")
	if nextChoice == "" {
		nextChoice = "resolve"
	}

	updatedChain, ok := k.state.updateActiveChain(chain.ID, chain.Version, func(state *ChainState) {
		state.StepCount++
		state.InferOutput = inferOutput
		state.CurrentPrimitive = nextChoice
		if nextChoice == "reference" {
			state.InferLoopCount++
		} else {
			state.InferLoopCount = 0
		}
	})
	if !ok {
		k.markStaleEvent(event, "chain changed during infer completion")
		return nil
	}
	if updatedChain.InferLoopCount >= maxReferenceInferHops {
		nextChoice = "resolve"
		updatedChain, _ = k.state.updateActiveChain(chain.ID, chain.Version, func(state *ChainState) {
			state.CurrentPrimitive = "resolve"
			state.InferLoopCount = 0
		})
		log.Printf(
			"kernel: infer loop limit reached chain=%s version=%d forcing=%s",
			updatedChain.ID,
			updatedChain.Version,
			nextChoice,
		)
	}
	log.Printf(
		"kernel: infer chose=%s chain=%s version=%d infer=%q",
		nextChoice,
		updatedChain.ID,
		updatedChain.Version,
		summarizeForLog(inferOutput, 160),
	)
	k.completeActiveStep(nextChoice, updatedChain.Perception, "", "", currentUnderstandingText(updatedChain.Perception.Understanding))

	if reachedChainLimit(updatedChain) {
		k.state.completeActiveChain(updatedChain.ID, updatedChain.Version, time.Now().UTC())
		k.publishChainState(nil)
		return nil
	}
	if nextChoice == "reference" {
		return k.enqueuePrimitive(ctx, updatedChain, EventKindReference)
	}
	return k.enqueuePrimitive(ctx, updatedChain, EventKindResolve)
}

func (k *Kernel) processResolveEvent(ctx context.Context, event *Event) error {
	chain, ok := k.state.activeChainForEvent(event)
	if !ok {
		k.markStaleEvent(event, "active chain changed before resolve")
		return nil
	}

	raw, _, err := k.runStep(ctx, *chain, "resolve", k.mustRender("resolve", map[string]string{
		"stimulus":         chain.Perception.Stimulus.RawPayload,
		"perception":       mustJSON(chain.Perception),
		"reference_output": chain.ReferenceOutput,
		"infer_output":     chain.InferOutput,
	}))
	if err != nil {
		if errors.Is(err, ErrStepInterrupted) {
			return k.handleInterruptedChain(chain.ID, chain.Version)
		}
		return err
	}
	text := extractSection(raw, "understanding")
	if text == "" {
		text = strings.TrimSpace(raw)
	}
	now := time.Now().UTC()
	understanding := &Understanding{
		ID:        newID(),
		Text:      text,
		UpdatedAt: now,
	}

	updatedChain, ok := k.state.updateActiveChain(chain.ID, chain.Version, func(state *ChainState) {
		state.StepCount++
		state.ResolveCount++
		state.CurrentPrimitive = string(EventKindExperience)
		state.ReferenceOutput = ""
		state.InferOutput = ""
		state.InferLoopCount = 0
		state.Perception.Understanding = understanding
		state.Perception.UpdatedAt = now
	})
	if !ok {
		k.markStaleEvent(event, "chain changed during resolve completion")
		return nil
	}
	log.Printf(
		"kernel: resolve updated understanding chain=%s version=%d understanding=%q",
		updatedChain.ID,
		updatedChain.Version,
		summarizeForLog(text, 180),
	)
	k.publishPerception(updatedChain.Perception)
	k.completeActiveStep("experience", updatedChain.Perception, "", "", currentUnderstandingText(understanding))

	if reachedChainLimit(updatedChain) || updatedChain.ResolveCount >= maxResolveCycles {
		k.state.completeActiveChain(updatedChain.ID, updatedChain.Version, time.Now().UTC())
		k.publishChainState(nil)
		return nil
	}
	return k.enqueuePrimitive(ctx, updatedChain, EventKindExperience)
}

func (k *Kernel) processInteractEvent(ctx context.Context, event *Event) error {
	chain, ok := k.state.activeChainForEvent(event)
	if !ok {
		k.markStaleEvent(event, "active chain changed before interact")
		return nil
	}

	understandingText := currentUnderstandingText(chain.Perception.Understanding)
	if understandingText == "" {
		if current := k.state.currentUnderstanding(); current != nil {
			understandingText = current.Text
		}
	}

	raw, _, err := k.runStep(ctx, *chain, "interact", k.mustRender("interact", map[string]string{
		"stimulus":      chain.Perception.Stimulus.RawPayload,
		"perception":    mustJSON(chain.Perception),
		"understanding": understandingText,
	}))
	if err != nil {
		if errors.Is(err, ErrStepInterrupted) {
			return k.handleInterruptedChain(chain.ID, chain.Version)
		}
		return err
	}

	channel := extractSection(raw, "interaction channel")
	if channel == "" {
		channel = "human"
	}
	message := extractSection(raw, "interaction message")
	if message == "" {
		message = strings.TrimSpace(raw)
	}

	now := time.Now().UTC()
	assistantMessage := k.state.appendInteraction("assistant", message, chain.RootStimulusID, now)
	k.state.appendHistory("assistant", "skyra", message, now)
	log.Printf(
		"kernel: interact emitted channel=%s chain=%s stimulus=%s message=%q",
		channel,
		chain.ID,
		chain.RootStimulusID,
		summarizeForLog(message, 180),
	)
	k.chatBroker.Publish(UIEvent{Type: "interaction.message", Data: assistantMessage})
	updatedChain, ok := k.state.updateActiveChain(chain.ID, chain.Version, func(state *ChainState) {
		state.StepCount++
		state.CurrentPrimitive = ""
	})
	if !ok {
		k.markStaleEvent(event, "chain changed during interact completion")
		return nil
	}
	k.completeActiveStep("complete", updatedChain.Perception, channel, message, understandingText)
	k.state.completeActiveChain(chain.ID, chain.Version, now)
	if updatedChain.ReviewingUser && k.state.hasSuspendedChains() {
		return k.resumeSuspendedChain(ctx)
	}
	k.publishChainState(nil)
	return nil
}

func (k *Kernel) enqueuePrimitive(ctx context.Context, chain *ChainState, kind EventKind) error {
	if chain == nil {
		return nil
	}
	event := &Event{
		ID:           newID(),
		StimulusID:   chain.RootStimulusID,
		Source:       chain.StimulusSource,
		Kind:         kind,
		Priority:     primitivePriority,
		ChainID:      chain.ID,
		ChainVersion: chain.Version,
		CreatedAt:    time.Now().UTC(),
	}
	if err := k.heap.Push(ctx, event); err != nil {
		return err
	}
	log.Printf(
		"kernel: enqueued event=%s kind=%s priority=%d chain=%s version=%d stimulus=%s queue_depth=%d",
		event.ID,
		event.Kind,
		event.Priority,
		event.ChainID,
		event.ChainVersion,
		event.StimulusID,
		k.heap.Len(),
	)
	return nil
}

func (k *Kernel) runStep(ctx context.Context, chain ChainState, frame string, prompt string) (string, ThoughtStepState, error) {
	startedAt := time.Now().UTC()
	step, ok := k.state.beginStep(chain.ID, chain.Version, frame, startedAt)
	if !ok {
		return "", ThoughtStepState{}, fmt.Errorf("chain became stale before step %q", frame)
	}
	log.Printf(
		"kernel: step started chain=%s version=%d step=%s index=%d frame=%s stimulus=%s content=%q",
		chain.ID,
		chain.Version,
		step.StepID,
		step.StepIndex,
		frame,
		chain.RootStimulusID,
		summarizeForLog(chain.Perception.Stimulus.RawPayload, 120),
	)
	k.thoughtBroker.Publish(UIEvent{Type: "thought.step.started", Data: step})

	stepCtx, cancel := context.WithCancel(ctx)
	k.setActiveStepCancel(cancel)
	defer func() {
		k.clearActiveStepCancel()
		cancel()
	}()

	raw, err := k.gateway.RunPrompt(stepCtx, prompt, func(delta string) {
		active := k.state.appendStepDelta(step.StepID, delta)
		if active != nil {
			k.thoughtBroker.Publish(UIEvent{Type: "thought.step.delta", Data: active})
		}
	})
	if err != nil {
		if active, ok := k.state.activeStepSnapshot(step.StepID); ok && active != nil {
			step = *active
		}
		finishedAt := time.Now().UTC()
		status := "error"
		if errors.Is(err, context.Canceled) {
			status = "interrupted"
		} else {
			step.RawOutput = "error: " + err.Error()
		}
		record := k.state.completeStep(step, status, "", chain.Perception, "", "", currentUnderstandingText(chain.Perception.Understanding), finishedAt)
		log.Printf(
			"kernel: step finished chain=%s step=%s index=%d frame=%s status=%s duration=%s raw=%q",
			chain.ID,
			step.StepID,
			step.StepIndex,
			frame,
			status,
			finishedAt.Sub(startedAt).Round(time.Millisecond),
			summarizeForLog(step.RawOutput, 200),
		)
		k.thoughtBroker.Publish(UIEvent{Type: "thought.step.completed", Data: record})
		if status == "interrupted" {
			return step.RawOutput, step, ErrStepInterrupted
		}
		return step.RawOutput, step, err
	}
	step.RawOutput = raw
	log.Printf(
		"kernel: step response chain=%s step=%s index=%d frame=%s duration=%s raw=%q",
		chain.ID,
		step.StepID,
		step.StepIndex,
		frame,
		time.Since(startedAt).Round(time.Millisecond),
		summarizeForLog(raw, 200),
	)
	return raw, step, nil
}

func (k *Kernel) completeActiveStep(primitiveChoice string, perception Perception, interactionChannel, interactionMessage, currentUnderstanding string) {
	active := k.state.thoughtSnapshot(k.heap.Len()).ActiveStep
	if active == nil {
		return
	}
	record := k.state.completeStep(*active, "completed", primitiveChoice, perception, interactionChannel, interactionMessage, currentUnderstanding, time.Now().UTC())
	k.thoughtBroker.Publish(UIEvent{Type: "thought.step.completed", Data: record})
}

func (k *Kernel) publishPerception(perception Perception) {
	k.thoughtBroker.Publish(UIEvent{Type: "thought.perception", Data: perception})
}

func (k *Kernel) publishChain(chain ChainState) {
	k.publishChainState(&chain)
}

func (k *Kernel) publishChainState(chain *ChainState) {
	if chain == nil {
		k.thoughtBroker.Publish(UIEvent{Type: "thought.chain", Data: map[string]any{"status": "idle"}})
		return
	}
	k.thoughtBroker.Publish(UIEvent{Type: "thought.chain", Data: chain})
}

func (k *Kernel) publishSuspendedChain(chain *ChainState) {
	if chain == nil {
		return
	}
	k.thoughtBroker.Publish(UIEvent{Type: "thought.chain.suspended", Data: chain})
}

func (k *Kernel) markStaleEvent(event *Event, reason string) {
	if event == nil {
		return
	}
	log.Printf(
		"kernel: marked stale event=%s kind=%s chain=%s version=%d stimulus=%s reason=%q",
		event.ID,
		event.Kind,
		emptyDash(event.ChainID),
		event.ChainVersion,
		emptyDash(event.StimulusID),
		reason,
	)
	k.thoughtBroker.Publish(UIEvent{Type: "thought.event.stale", Data: map[string]any{
		"event_id":      event.ID,
		"kind":          event.Kind,
		"chain_id":      event.ChainID,
		"chain_version": event.ChainVersion,
		"stimulus_id":   event.StimulusID,
		"stale_reason":  reason,
		"marked_at":     time.Now().UTC(),
	}})
}

func (k *Kernel) handleInterruptedChain(chainID string, version int) error {
	suspended, ok := k.state.suspendActiveChain(chainID, version, time.Now().UTC())
	if ok {
		log.Printf(
			"kernel: interrupted active chain=%s version=%d current_primitive=%s",
			suspended.ID,
			suspended.Version,
			emptyDash(suspended.CurrentPrimitive),
		)
		k.publishSuspendedChain(suspended)
		k.publishChainState(nil)
	}
	return ErrStepInterrupted
}

func (k *Kernel) resumeSuspendedChain(ctx context.Context) error {
	resumed, ok := k.state.resumeSuspendedChain(time.Now().UTC())
	if !ok || resumed == nil {
		k.publishChainState(nil)
		return nil
	}
	log.Printf(
		"kernel: resumed suspended chain=%s version=%d at=%s stimulus=%s",
		resumed.ID,
		resumed.Version,
		EventKindReference,
		resumed.RootStimulusID,
	)
	k.publishPerception(resumed.Perception)
	k.publishChain(*resumed)
	return k.enqueuePrimitive(ctx, resumed, EventKindReference)
}

func (k *Kernel) interruptActiveStep() {
	k.stepMu.Lock()
	cancel := k.stepCancel
	k.stepMu.Unlock()
	if cancel != nil {
		log.Printf("kernel: canceling active model step")
		cancel()
	}
}

func (k *Kernel) setActiveStepCancel(cancel context.CancelFunc) {
	k.stepMu.Lock()
	defer k.stepMu.Unlock()
	k.stepCancel = cancel
}

func (k *Kernel) clearActiveStepCancel() {
	k.stepMu.Lock()
	defer k.stepMu.Unlock()
	k.stepCancel = nil
}

func (k *Kernel) lookupStimulus(id string) (Stimulus, bool) {
	k.state.mu.RLock()
	defer k.state.mu.RUnlock()
	for i := len(k.state.stimuli) - 1; i >= 0; i-- {
		if k.state.stimuli[i].ID == id {
			return k.state.stimuli[i], true
		}
	}
	return Stimulus{}, false
}

func (k *Kernel) mustRender(name string, values map[string]string) string {
	rendered, err := k.prompts.Render(name, values)
	if err != nil {
		panic(err)
	}
	return rendered
}

func extractSection(raw, label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		return ""
	}
	pattern := fmt.Sprintf(`(?is)%s:\s*(.*?)\s*(?:\n[a-z][a-z _-]*:\s*|\z)`, regexp.QuoteMeta(label))
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(strings.TrimSpace(raw))
	if len(match) < 2 {
		return ""
	}
	return strings.TrimSpace(match[1])
}

func parseChoice(raw string, allowed ...string) string {
	candidate := strings.ToLower(strings.TrimSpace(raw))
	if candidate == "" {
		return ""
	}
	for _, allowedChoice := range allowed {
		if candidate == allowedChoice {
			return allowedChoice
		}
	}
	for _, allowedChoice := range allowed {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(allowedChoice) + `\b`)
		if re.FindStringIndex(candidate) != nil {
			return allowedChoice
		}
	}
	return ""
}

func currentUnderstandingText(value *Understanding) string {
	if value == nil {
		return ""
	}
	return value.Text
}

func reachedChainLimit(chain *ChainState) bool {
	return chain != nil && chain.StepCount >= maxChainStepCount
}

func summarizeForLog(value string, limit int) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = strings.Join(strings.Fields(value), " ")
	if limit <= 0 || len(value) <= limit {
		return value
	}
	if limit <= 3 {
		return value[:limit]
	}
	return value[:limit-3] + "..."
}

func emptyDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func isHumanSource(source string) bool {
	source = strings.TrimSpace(strings.ToLower(source))
	return source == "human" || source == "user"
}
