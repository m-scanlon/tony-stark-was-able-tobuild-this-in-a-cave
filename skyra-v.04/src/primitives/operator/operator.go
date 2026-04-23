package operator

import "skyra-v04/src/primitives/entity"

// IOperator is the marker interface for verbs beings can speak during conversation.
// Operators are entities that route and act — they have no medium of their own.
// They live in the world's EntityMap and are listed in each being's ~operators.
type IOperator interface {
	entity.Entity
}
