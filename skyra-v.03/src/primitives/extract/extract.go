package extract

import (
	"fmt"
	"strings"
)

func Meaning(expression, token, name string) (string, error) {
	idx := strings.Index(expression, token)
	if idx == -1 {
		return "", fmt.Errorf("%s: token %q not found in expression", name, token)
	}

	rest := strings.TrimSpace(expression[idx+len(token):])
	if rest == "" {
		return "", fmt.Errorf("%s: no value after token %q", name, token)
	}

	// value ends at the next ~ token or | divider
	end := len(rest)
	for _, delim := range []string{"~", "|"} {
		if i := strings.Index(rest, delim); i != -1 && i < end {
			end = i
		}
	}

	value := strings.TrimSpace(rest[:end])
	if value == "" {
		return "", fmt.Errorf("%s: empty value for token %q", name, token)
	}
	return value, nil
}
