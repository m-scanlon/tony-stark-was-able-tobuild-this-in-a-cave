package metaxu

import (
	"fmt"
	"strings"

	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/meaning"
	"skyra-v03/src/world"
)

type Signal struct {
	ID             string
	Origin         string
	TraceToken     string
	ThreadID       string
	About          string
	Because        string
	ContextEntries []being.ExchangeEntry
	Impulse        string
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
	ThreadID          string
	WrittenBeingName  string
	WrittenPeerName   string
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
		ThreadID:      signal.ThreadID,
	}

	target, ok := m.world.BeingByName(parsed.TargetName)
	if !ok {
		result.DropReason = "target could not be resolved"
		return result
	}
	result.TargetName = target.Name
	result.ReceiverCognitive = target.Cognitive

	operatorName := exchangeOperator(parsed.Expression)
	if operatorName != "" && target.Name != operatorName {
		result.DropReason = fmt.Sprintf("target %s directly; do not wrap %s inside a signal to %s", operatorName, operatorName, target.Name)
		return result
	}

	threadID := signal.ThreadID
	if threadID == "" && origin.Cognitive && target.Cognitive && operatorName == "" {
		if id, ok := m.world.FindOpenExchangeThread(origin.Name, target.Name); ok {
			threadID = id
		}
	}
	result.ThreadID = threadID

	if origin.Cognitive && target.Cognitive && operatorName == "" && !m.world.HasExchangeThread(origin.Name, target.Name, threadID) {
		result.DropReason = fmt.Sprintf("no open exchange with %s for thread %s; open one with start-exchange first", target.Name, threadID)
		return result
	}

	// Extract ~expression-reference from expression if present; resolve context entries and strip from raw
	ref, _ := meaning.Extract(parsed.Expression, "~expression-reference", "continue")
	rawImpulse := being.Impulse(parsed.Raw)
	var refContextEntries []being.ExchangeEntry
	if ref != "" {
		refContextEntries = m.world.ResolveExpressionRef(origin.Name, threadID, ref)
		rawImpulse = being.Impulse(meaning.Strip(parsed.Raw, "~expression-reference"))
	}

	contextEntries := signal.ContextEntries
	if len(refContextEntries) > 0 {
		contextEntries = append(contextEntries, refContextEntries...)
	}

	delivery := being.DeliveredImpulse{
		OriginName:     origin.Name,
		ThreadID:       threadID,
		About:          signal.About,
		Because:        signal.Because,
		ContextEntries: contextEntries,
		Raw:            rawImpulse,
		Parsed:         parsed,
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
	} else {
		if origin.Cognitive && operatorName == "" {
			if _, err := origin.SendToPeer(target.Name, delivery); err != nil {
				result.DropReason = err.Error()
				return result
			}
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
	}

	result.Status = RouteStatusRouted
	result.WrittenBeingName = target.Name
	result.WrittenPeerName = origin.Name

	if present, err := target.DerivePresent(origin); err == nil {
		result.ReceiverPresent = present
	}
	return result
}

func exchangeOperator(expression string) string {
	expression = strings.TrimSpace(expression)
	for _, operatorName := range []string{"start-exchange", "close-exchange"} {
		if expression == operatorName || strings.HasPrefix(expression, operatorName+" ") {
			return operatorName
		}
	}
	return ""
}
