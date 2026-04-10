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
	Status           RouteStatus
	DropReason       string
	ParsedImpulse    *being.ParsedImpulse
	OriginName       string
	TargetName       string
	WrittenBeingName string
	WrittenPeerName  string
	NewExchange      bool
	ReceiverPresent  string
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

	source, ok := m.world.BeingByName(parsed.Source)
	if !ok {
		result.DropReason = "source could not be resolved"
		return result
	}

	target, ok := m.world.BeingByName(parsed.TargetName)
	if !ok {
		result.DropReason = "target could not be resolved"
		return result
	}
	result.TargetName = target.Name

	originDelivery := being.DeliveredImpulse{
		Raw:    being.Impulse(parsed.Raw),
		Parsed: parsed,
	}

	targetDelivery := being.DeliveredImpulse{
		OriginName: origin.Name,
		Raw:        being.Impulse(parsed.Raw),
		Parsed:     parsed,
	}

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

	if _, err := origin.SendToPeer(source.Name, originDelivery); err != nil {
		result.DropReason = err.Error()
		return result
	}

	if source.Name != target.Name {
		if _, err := origin.SendToPeer(target.Name, originDelivery); err != nil {
			result.DropReason = err.Error()
			return result
		}
	}

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
