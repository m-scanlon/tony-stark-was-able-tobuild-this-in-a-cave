package kernel

import (
	"testing"

	"skyra-v03/src/domain"
)

func TestAcceptSignalRoutesToReceiverAndStartsNewExchange(t *testing.T) {
	state, origin, target := seededKernelState(t)
	k := New(state)

	result := k.AcceptSignal(domain.Signal{
		ID:         "sig-1",
		Origin:     origin.Name,
		TraceToken: "trace-1",
		Impulse:    "skyra skyra hello there | skyra: greeting",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if !result.NewExchange {
		t.Fatalf("NewExchange = false, want true")
	}
	if result.WrittenBeingName != target.Name {
		t.Fatalf("WrittenBeingName = %q, want %q", result.WrittenBeingName, target.Name)
	}
	if result.WrittenPeerName != origin.Name {
		t.Fatalf("WrittenPeerName = %q, want %q", result.WrittenPeerName, origin.Name)
	}
	want := "name: skyra\nidentity: system\npurpose: relate\n\nYou are in an exchange with: michael\nthe identity of michael is: builder\nthe purpose of michael is: hold the line\n\nmichael: hello there\n\nrelationships:\nCall any of your relationships using this syntax-\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>\n<being> must be one of your relationships listed below\n<source> is the being you are currently in exchange with\n<reason> is why you are firing this expression\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nbuilder - hold the line"
	if result.ReceiverPresent != want {
		t.Fatalf("ReceiverPresent = %q, want %q", result.ReceiverPresent, want)
	}
}

func TestAcceptSignalAppendsToExistingOpenExchange(t *testing.T) {
	state, origin, target := seededKernelState(t)
	k := New(state)

	first, err := domain.NewImpulse("skyra skyra first | skyra: opening")
	if err != nil {
		t.Fatalf("NewImpulse(first) error = %v", err)
	}
	if _, err := target.SendToPeer(origin.Name, deliveredFrom(t, origin.Name, first)); err != nil {
		t.Fatalf("SendToPeer(first) error = %v", err)
	}

	result := k.AcceptSignal(domain.Signal{
		ID:         "sig-2",
		Origin:     origin.Name,
		TraceToken: "trace-2",
		Impulse:    "skyra skyra second | skyra: follow-up",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if result.NewExchange {
		t.Fatalf("NewExchange = true, want false")
	}
	want := "name: skyra\nidentity: system\npurpose: relate\n\nYou are in an exchange with: michael\nthe identity of michael is: builder\nthe purpose of michael is: hold the line\n\nmichael: first\n\nmichael: second\n\nrelationships:\nCall any of your relationships using this syntax-\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>\n<being> must be one of your relationships listed below\n<source> is the being you are currently in exchange with\n<reason> is why you are firing this expression\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nbuilder - hold the line"
	if result.ReceiverPresent != want {
		t.Fatalf("ReceiverPresent = %q, want %q", result.ReceiverPresent, want)
	}
}

func TestAcceptSignalStartsNewExchangeAfterClosedTop(t *testing.T) {
	state, origin, target := seededKernelState(t)
	k := New(state)

	first, err := domain.NewImpulse("skyra skyra first | skyra: opening")
	if err != nil {
		t.Fatalf("NewImpulse(first) error = %v", err)
	}
	closeImpulse, err := domain.NewImpulse("skyra skyra -close | skyra: closing")
	if err != nil {
		t.Fatalf("NewImpulse(close) error = %v", err)
	}
	if _, err := target.SendToPeer(origin.Name, deliveredFrom(t, origin.Name, first)); err != nil {
		t.Fatalf("SendToPeer(first) error = %v", err)
	}
	if _, err := target.SendToPeer(origin.Name, deliveredFrom(t, origin.Name, closeImpulse)); err != nil {
		t.Fatalf("SendToPeer(close) error = %v", err)
	}

	result := k.AcceptSignal(domain.Signal{
		ID:         "sig-3",
		Origin:     origin.Name,
		TraceToken: "trace-3",
		Impulse:    "skyra skyra after close | skyra: resuming",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if !result.NewExchange {
		t.Fatalf("NewExchange = false, want true")
	}

	targetPeer, ok := target.PeerByName(origin.Name)
	if !ok {
		t.Fatalf("PeerByName() ok = false, want true")
	}
	stack, ok := targetPeer.(*domain.ExchangeStack)
	if !ok {
		t.Fatalf("target peer type = %T, want *domain.ExchangeStack", targetPeer)
	}
	if len(stack.Exchanges()) != 2 {
		t.Fatalf("len(targetPeer.Exchanges()) = %d, want 2", len(stack.Exchanges()))
	}
	if len(stack.Exchanges()[0]) != 2 {
		t.Fatalf("len(targetPeer.Exchanges()[0]) = %d, want 2", len(stack.Exchanges()[0]))
	}
	if len(stack.Exchanges()[1]) != 1 {
		t.Fatalf("len(targetPeer.Exchanges()[1]) = %d, want 1", len(stack.Exchanges()[1]))
	}
	want := "name: skyra\nidentity: system\npurpose: relate\n\nYou are in an exchange with: michael\nthe identity of michael is: builder\nthe purpose of michael is: hold the line\n\nmichael: after close\n\nrelationships:\nCall any of your relationships using this syntax-\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>\n<being> must be one of your relationships listed below\n<source> is the being you are currently in exchange with\n<reason> is why you are firing this expression\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nbuilder - hold the line"
	if result.ReceiverPresent != want {
		t.Fatalf("ReceiverPresent = %q, want %q", result.ReceiverPresent, want)
	}
}

func TestAcceptSignalCloseWritesOnlyToOriginStack(t *testing.T) {
	state, origin, target := seededKernelState(t)
	k := New(state)

	open, err := domain.NewImpulse("skyra skyra open | skyra: opening")
	if err != nil {
		t.Fatalf("NewImpulse(open) error = %v", err)
	}
	if _, err := origin.EmitToPeer(target.Name, open); err != nil {
		t.Fatalf("EmitToPeer(open) error = %v", err)
	}

	result := k.AcceptSignal(domain.Signal{
		ID:         "sig-4",
		Origin:     origin.Name,
		TraceToken: "trace-4",
		Impulse:    "skyra skyra -close | skyra: closing",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if result.WrittenBeingName != origin.Name {
		t.Fatalf("WrittenBeingName = %q, want %q", result.WrittenBeingName, origin.Name)
	}
	if result.WrittenPeerName != target.Name {
		t.Fatalf("WrittenPeerName = %q, want %q", result.WrittenPeerName, target.Name)
	}

	originPeer, ok := origin.PeerByName(target.Name)
	if !ok {
		t.Fatalf("origin PeerByName() ok = false, want true")
	}
	targetPeer, ok := target.PeerByName(origin.Name)
	if !ok {
		t.Fatalf("target PeerByName() ok = false, want true")
	}
	originStack, ok := originPeer.(*domain.ExchangeStack)
	if !ok {
		t.Fatalf("origin peer type = %T, want *domain.ExchangeStack", originPeer)
	}
	targetStack, ok := targetPeer.(*domain.ExchangeStack)
	if !ok {
		t.Fatalf("target peer type = %T, want *domain.ExchangeStack", targetPeer)
	}
	if len(originStack.Exchanges()) != 1 {
		t.Fatalf("len(originStack.Exchanges()) = %d, want 1", len(originStack.Exchanges()))
	}
	if len(originStack.Exchanges()[0]) != 2 {
		t.Fatalf("len(originStack.Exchanges()[0]) = %d, want 2", len(originStack.Exchanges()[0]))
	}
	if len(targetStack.Exchanges()) != 0 {
		t.Fatalf("len(targetStack.Exchanges()) = %d, want 0", len(targetStack.Exchanges()))
	}
	if string(originStack.Exchanges()[0][1]) != "skyra skyra -close | skyra: closing" {
		t.Fatalf("close stored raw impulse = %q, want raw value", string(originStack.Exchanges()[0][1]))
	}
}

func TestAcceptSignalDropsWhenTargetCannotBeResolved(t *testing.T) {
	state, origin, _ := seededKernelState(t)
	k := New(state)

	result := k.AcceptSignal(domain.Signal{
		ID:         "sig-5",
		Origin:     origin.Name,
		TraceToken: "trace-5",
		Impulse:    "skyra unknown hello | skyra: testing",
	})

	if result.Status != RouteStatusDropped {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusDropped)
	}
	if result.DropReason == "" {
		t.Fatalf("DropReason = empty, want value")
	}
}

func TestAcceptSignalRoutesToExternalDispatchAndDerivesExpressionOnlyPresent(t *testing.T) {
	state := domain.NewKernelState()
	origin, err := domain.NewBeing(
		"michael",
		domain.Nature{Identity: "builder", Purpose: "hold the line"},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(origin) error = %v", err)
	}
	target, err := domain.NewBeing(
		"sensor",
		domain.Nature{Identity: "sensor", Purpose: "receive"},
		true,
		false,
	)
	if err != nil {
		t.Fatalf("NewBeing(target) error = %v", err)
	}
	if err := state.InsertBeing(origin); err != nil {
		t.Fatalf("InsertBeing(origin) error = %v", err)
	}
	if err := state.InsertBeing(target); err != nil {
		t.Fatalf("InsertBeing(target) error = %v", err)
	}
	if err := origin.SeedPeer(target.Name, target.Nature); err != nil {
		t.Fatalf("SeedPeer(origin->target) error = %v", err)
	}
	dispatch, err := domain.NewExternalDispatch(origin.Name, origin.Nature)
	if err != nil {
		t.Fatalf("NewExternalDispatch() error = %v", err)
	}
	if err := target.AttachPeer(dispatch); err != nil {
		t.Fatalf("AttachPeer(target<-origin) error = %v", err)
	}

	k := New(state)
	result := k.AcceptSignal(domain.Signal{
		ID:         "sig-6",
		Origin:     origin.Name,
		TraceToken: "trace-6",
		Impulse:    "skyra sensor heat spike | sensor: reporting",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if result.ReceiverPresent != "heat spike" {
		t.Fatalf("ReceiverPresent = %q, want %q", result.ReceiverPresent, "heat spike")
	}
}

func seededKernelState(t *testing.T) (*domain.KernelState, *domain.Being, *domain.Being) {
	t.Helper()

	state := domain.NewKernelState()

	origin, err := domain.NewBeing(
		"michael",
		domain.Nature{Identity: "builder", Purpose: "hold the line"},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(origin) error = %v", err)
	}

	target, err := domain.NewBeing(
		"skyra",
		domain.Nature{Identity: "system", Purpose: "relate"},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(target) error = %v", err)
	}

	if err := state.InsertBeing(origin); err != nil {
		t.Fatalf("InsertBeing(origin) error = %v", err)
	}
	if err := state.InsertBeing(target); err != nil {
		t.Fatalf("InsertBeing(target) error = %v", err)
	}
	if err := state.SeedRelationship(origin.Name, target.Name); err != nil {
		t.Fatalf("SeedRelationship() error = %v", err)
	}

	return state, origin, target
}

func deliveredFrom(t *testing.T, originName string, impulse domain.Impulse) domain.DeliveredImpulse {
	t.Helper()

	parsed, err := impulse.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	return domain.DeliveredImpulse{
		OriginName: originName,
		Raw:        impulse,
		Parsed:     parsed,
	}
}
