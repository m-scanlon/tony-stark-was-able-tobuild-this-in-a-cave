package kernel

import (
	"skyra-v03/src/domain"
)

type RouteStatus string

const (
	RouteStatusRouted  RouteStatus = "routed"
	RouteStatusDropped RouteStatus = "dropped"
)

type Result struct {
	Status          RouteStatus
	DropReason      string
	ParsedImpulse   *domain.ParsedImpulse
	OriginName      string
	TargetName      string
	WrittenBeingName string
	WrittenPeerName  string
	NewExchange     bool
	ReceiverPresent string
}

type Kernel struct {
	state *domain.KernelState
}

func New(state *domain.KernelState) *Kernel {
	return &Kernel{state: state}
}

func (k *Kernel) AcceptSignal(signal domain.Signal) Result {
	if k == nil || k.state == nil {
		return Result{
			Status:     RouteStatusDropped,
			DropReason: "kernel state is unavailable",
		}
	}

	// Resolve origin — hard drop if not found
	origin, ok := k.state.BeingByName(signal.Origin)
	if !ok {
		return Result{
			Status:     RouteStatusDropped,
			DropReason: "origin could not be resolved",
		}
	}

	parsed, err := domain.ParseImpulse(signal.Impulse)
	if err != nil {
		return Result{
			Status:     RouteStatusDropped,
			DropReason: err.Error(),
			OriginName: origin.Name,
		}
	}

	result := Result{
		Status:        RouteStatusDropped,
		ParsedImpulse: &parsed,
		OriginName:    origin.Name,
	}

	// Resolve source — silent drop if not found
	source, ok := k.state.BeingByName(parsed.Source)
	if !ok {
		result.DropReason = "source could not be resolved"
		return result
	}

	// Resolve target — bounce to origin's exchange with source if not found
	target, ok := k.state.BeingByName(parsed.TargetName)
	if !ok {
		// TODO: write error impulse to origin's exchange with source
		result.DropReason = "target could not be resolved"
		return result
	}
	result.TargetName = target.Name

	// Delivery for origin-side writes (no OriginName — prevents target swap)
	originDelivery := domain.DeliveredImpulse{
		Raw:    domain.Impulse(parsed.Raw),
		Parsed: parsed,
	}

	// Delivery for target-side write (with OriginName — triggers target swap)
	targetDelivery := domain.DeliveredImpulse{
		OriginName: origin.Name,
		Raw:        domain.Impulse(parsed.Raw),
		Parsed:     parsed,
	}

	// Close: write only to origin's exchange with target
	if parsed.IsClose() {
		channelResult, err := origin.SendToPeer(target.Name, originDelivery)
		if err != nil {
			result.DropReason = err.Error()
			return result
		}
		if !channelResult.Routed {
			result.DropReason = channelResult.DropReason
			return result
		}

		result.Status = RouteStatusRouted
		result.NewExchange = channelResult.NewExchange
		result.WrittenBeingName = origin.Name
		result.WrittenPeerName = target.Name

		if present, err := target.DerivePresent(origin); err == nil {
			result.ReceiverPresent = present
		}
		return result
	}

	// Write 1: origin's exchange with source
	if _, err := origin.SendToPeer(source.Name, originDelivery); err != nil {
		result.DropReason = err.Error()
		return result
	}

	// Write 2: origin's exchange with target (skip if source == target — same exchange)
	if source.Name != target.Name {
		if _, err := origin.SendToPeer(target.Name, originDelivery); err != nil {
			result.DropReason = err.Error()
			return result
		}
	}

	// Write 3: target's exchange with origin
	channelResult, err := target.SendToPeer(origin.Name, targetDelivery)
	if err != nil {
		result.DropReason = err.Error()
		return result
	}
	if !channelResult.Routed {
		result.DropReason = channelResult.DropReason
		return result
	}

	result.Status = RouteStatusRouted
	result.NewExchange = channelResult.NewExchange
	result.WrittenBeingName = target.Name
	result.WrittenPeerName = origin.Name

	if present, err := target.DerivePresent(origin); err == nil {
		result.ReceiverPresent = present
	}
	return result
}
