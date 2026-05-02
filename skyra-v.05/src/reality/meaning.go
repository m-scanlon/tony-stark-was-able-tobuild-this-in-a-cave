package reality

import (
	"fmt"
	"strings"
)

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

func ExtractTag(text, tag string) (string, error) {
	open := "<" + tag + ">"
	close := "</" + tag + ">"
	start := strings.Index(text, open)
	if start == -1 {
		return "", fmt.Errorf("tag %q not found", tag)
	}
	after := text[start+len(open):]
	end := strings.Index(after, close)
	if end == -1 {
		return "", fmt.Errorf("tag %q not closed", tag)
	}
	return strings.TrimSpace(after[:end]), nil
}

func StripTag(text, tag string) string {
	open := "<" + tag + ">"
	close := "</" + tag + ">"
	start := strings.Index(text, open)
	if start == -1 {
		return text
	}
	after := text[start+len(open):]
	end := strings.Index(after, close)
	if end == -1 {
		return text
	}
	before := strings.TrimSpace(text[:start])
	rest := strings.TrimSpace(after[end+len(close):])
	if before == "" {
		return rest
	}
	if rest == "" {
		return before
	}
	return before + " " + rest
}

