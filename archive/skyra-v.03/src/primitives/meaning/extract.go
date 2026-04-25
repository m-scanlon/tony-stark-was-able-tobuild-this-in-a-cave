package meaning

import (
	"fmt"
	"strings"
)

// ExtractToEnd extracts the value after token, stopping only at | (not at ~).
// Use this when the value itself may contain ~ flags, e.g. expression syntax.
func ExtractToEnd(expression, token, name string) (string, error) {
	idx := strings.Index(expression, token)
	if idx == -1 {
		return "", fmt.Errorf("%s: token %q not found in expression", name, token)
	}

	rest := strings.TrimSpace(expression[idx+len(token):])
	if rest == "" {
		return "", fmt.Errorf("%s: no value after token %q", name, token)
	}

	end := len(rest)
	if i := strings.Index(rest, "|"); i != -1 && i < end {
		end = i
	}

	value := strings.TrimSpace(rest[:end])
	if value == "" {
		return "", fmt.Errorf("%s: empty value for token %q", name, token)
	}
	return value, nil
}

func Extract(expression, token, name string) (string, error) {
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
