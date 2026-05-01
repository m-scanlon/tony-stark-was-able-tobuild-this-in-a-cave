package medium

import (
	"skyra-v04/src/inference"
	"skyra-v04/src/primitives/entity"
)

func init() {
	Register("inference", func(present string, _ entity.Relation) (string, error) {
		return inference.Call(present)
	})
}
