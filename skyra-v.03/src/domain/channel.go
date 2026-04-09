package domain

type DeliveredImpulse struct {
	OriginName string
	Raw        Impulse
	Parsed     ParsedImpulse
}

type ChannelResult struct {
	Routed      bool
	NewExchange bool
	DropReason  string
}

type RelationshipChannel interface {
	Send(delivery DeliveredImpulse) ChannelResult
	Name() string
	PeerNature() Nature
}

type presentDeriver interface {
	derivePresent(receiver *Being, sender *Being) string
}
