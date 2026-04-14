package metaxu

import (
	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/world"
)

type Signal struct {
	ID         string
	Origin     string
	TraceToken string
	Impulse    string
}

type RouteStatus string

const (
	RouteStatusRouted  RouteStatus = "routed"
	RouteStatusDropped RouteStatus = "dropped"
)

type Result struct {
	Status            RouteStatus
	DropReason        string
	ParsedImpulse     *being.ParsedImpulse
	OriginName        string
	TargetName        string
	WrittenBeingName  string
	WrittenPeerName   string
	NewExchange       bool
	ReceiverPresent   string
	ReceiverCognitive bool
}

type Metaxu struct {
	world *world.World
}

func New(w *world.World) *Metaxu {
	return &Metaxu{world: w}
}

func (m *Metaxu) AcceptSignal(signal Signal) Result {
	if m == nil || m.world == nil {
		return Result{
			Status:     RouteStatusDropped,
			DropReason: "metaxu world is unavailable",
		}
	}

	origin, ok := m.world.BeingByName(signal.Origin)
	if !ok {
		return Result{
			Status:     RouteStatusDropped,
			DropReason: "origin could not be resolved",
		}
	}

	parsed, err := being.ParseImpulse(signal.Impulse)
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

	target, ok := m.world.BeingByName(parsed.TargetName)
	if !ok {
		result.DropReason = "target could not be resolved"
		return result
	}
	result.TargetName = target.Name
	result.ReceiverCognitive = target.Cognitive

	delivery := being.DeliveredImpulse{
		OriginName: origin.Name,
		Raw:        being.Impulse(parsed.Raw),
		Parsed:     parsed,
	}

	if parsed.IsClose() {
		channelResult, err := origin.SendToPeer(target.Name, delivery)
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

	if origin.Name == target.Name {
		// self-call: write once to the self channel
		channelResult, err := origin.SendToPeer(origin.Name, delivery)
		if err != nil {
			result.DropReason = err.Error()
			return result
		}
		if !channelResult.Routed {
			result.DropReason = channelResult.DropReason
			return result
		}
		result.NewExchange = channelResult.NewExchange
	} else {
		if _, err := origin.SendToPeer(target.Name, delivery); err != nil {
			result.DropReason = err.Error()
			return result
		}
		channelResult, err := target.SendToPeer(origin.Name, delivery)
		if err != nil {
			result.DropReason = err.Error()
			return result
		}
		if !channelResult.Routed {
			result.DropReason = channelResult.DropReason
			return result
		}
		result.NewExchange = channelResult.NewExchange
	}

	result.Status = RouteStatusRouted
	result.WrittenBeingName = target.Name
	result.WrittenPeerName = origin.Name

	if present, err := target.DerivePresent(origin); err == nil {
		result.ReceiverPresent = present
	}
	return result
}
