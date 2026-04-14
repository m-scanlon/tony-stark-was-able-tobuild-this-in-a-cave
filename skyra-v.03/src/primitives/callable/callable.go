package callable

import "skyra-v03/src/primitives/extract"

type Callable struct {
	Value string
}

func CreateCallable(expression string) (Callable, error) {
	value, err := extract.Meaning(expression, "~callable", "callable")
	if err != nil {
		return Callable{}, nil // optional field
	}
	return Callable{Value: value}, nil
}
