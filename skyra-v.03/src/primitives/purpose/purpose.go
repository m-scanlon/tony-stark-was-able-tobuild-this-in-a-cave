package purpose

import "skyra-v03/src/primitives/extract"

type Purpose struct {
	Value string
}

func CreatePurpose(expression string) (Purpose, error) {
	value, err := extract.Meaning(expression, "~purpose", "purpose")
	if err != nil {
		return Purpose{}, err
	}
	return Purpose{Value: value}, nil
}
