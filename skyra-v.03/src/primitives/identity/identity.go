package identity

import "skyra-v03/src/primitives/extract"

type Identity struct {
	Value string
}

func CreateIdentity(expression string) (Identity, error) {
	value, err := extract.Meaning(expression, "~identity", "identity")
	if err != nil {
		return Identity{}, err
	}
	return Identity{Value: value}, nil
}
