package language

import "skyra-v03/src/primitives/meaning"

type Language struct {
	Value string
}

func CreateLanguage(expression string) (Language, error) {
	value, err := meaning.Extract(expression, "~expression", "language")
	if err != nil {
		return Language{}, err
	}
	return Language{Value: value}, nil
}
