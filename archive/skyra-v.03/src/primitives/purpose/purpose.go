package purpose

import "skyra-v03/src/primitives/meaning"

type Purpose struct {
	Value string
}

func CreatePurpose(expression string) (Purpose, error) {
	value, err := meaning.Extract(expression, "~purpose", "purpose")
	if err != nil {
		return Purpose{}, err
	}
	return Purpose{Value: value}, nil
}
