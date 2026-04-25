package identity

import "skyra-v03/src/primitives/meaning"

type Identity struct {
	Value string
}

func CreateIdentity(expression string) (Identity, error) {
	value, err := meaning.Extract(expression, "~identity", "identity")
	if err != nil {
		return Identity{}, err
	}
	return Identity{Value: value}, nil
}
