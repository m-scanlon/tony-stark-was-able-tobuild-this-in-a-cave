package meaning

import (
	"fmt"
	"strings"
)

// Extract reads the value after token. By default, the value ends at the next
// ~ flag or | divider. Pass delimiters to override — e.g. Extract(exp, token, name, "|")
// to read to the end of the expression (stopping only at |).
func Extract(expression, token, name string, delimiters ...string) (string, error) {
	idx := strings.Index(expression, token)
	if idx == -1 {
		return "", fmt.Errorf("%s: token %q not found in expression", name, token)
	}

	rest := strings.TrimSpace(expression[idx+len(token):])
	if rest == "" {
		return "", fmt.Errorf("%s: no value after token %q", name, token)
	}

	if len(delimiters) == 0 {
		delimiters = []string{"~", "|"}
	}

	end := len(rest)
	for _, delim := range delimiters {
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

// Strip removes ~token <value> from expression, where value ends at the next ~ or |.
func Strip(expression, token string) string {
	idx := strings.Index(expression, token)
	if idx == -1 {
		return expression
	}
	before := strings.TrimSpace(expression[:idx])
	rest := expression[idx+len(token):]
	end := len(rest)
	for _, delim := range []string{"~", "|"} {
		if i := strings.Index(rest, delim); i != -1 && i < end {
			end = i
		}
	}
	after := strings.TrimSpace(rest[end:])
	if before == "" {
		return after
	}
	if after == "" {
		return before
	}
	return before + " " + after
}
