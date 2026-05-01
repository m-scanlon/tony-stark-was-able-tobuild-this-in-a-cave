// Package skyraproto is the wire-format library for compiled beings.
// Skyra writes a binary that imports this package, reads a relation from stdin,
// does its work, and writes a relation (or many) back to stdout.
package skyraproto

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Relation struct {
	ID       string
	Origin   string
	ThreadID string
	Impulse  string
}

// ReadRelation reads one protocol line from r and parses it into a Relation.
// The protocol line has shape: skyra <operator> <args> | <reason>
func ReadRelation(r io.Reader) (Relation, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return Relation{}, err
		}
		return Relation{}, io.EOF
	}
	return Parse(strings.TrimSpace(scanner.Text()))
}

// Parse turns a raw protocol line into a Relation.
func Parse(raw string) (Relation, error) {
	if raw == "" {
		return Relation{}, fmt.Errorf("skyraproto: empty input")
	}
	left := raw
	if i := strings.Index(raw, "|"); i != -1 {
		left = strings.TrimSpace(raw[:i])
	}
	tokens := strings.Fields(left)
	if len(tokens) < 2 {
		return Relation{}, fmt.Errorf("skyraproto: expected at least protocol and target")
	}
	if tokens[0] != "skyra" {
		return Relation{}, fmt.Errorf("skyraproto: must begin with skyra")
	}
	return Relation{
		ID:      tokens[1],
		Impulse: strings.Join(tokens[2:], " "),
	}, nil
}

// WriteRelation emits a protocol line to w.
// Shape: skyra <target> <impulse> | <reason>
func WriteRelation(w io.Writer, target, impulse, reason string) error {
	_, err := fmt.Fprintf(w, "skyra %s %s | %s\n", target, impulse, reason)
	return err
}

// Extract pulls a ~flag value out of an impulse. By default, the value ends at
// the next ~ flag or | divider. Pass a single custom delimiter to override
// (e.g. only split at |).
func Extract(impulse, flag string, stopAt ...string) (string, bool) {
	idx := strings.Index(impulse, flag)
	if idx == -1 {
		return "", false
	}
	rest := strings.TrimSpace(impulse[idx+len(flag):])
	if rest == "" {
		return "", false
	}
	delims := []string{"~", "|"}
	if len(stopAt) > 0 {
		delims = stopAt
	}
	end := len(rest)
	for _, d := range delims {
		if i := strings.Index(rest, d); i != -1 && i < end {
			end = i
		}
	}
	value := strings.TrimSpace(rest[:end])
	if value == "" {
		return "", false
	}
	return value, true
}
