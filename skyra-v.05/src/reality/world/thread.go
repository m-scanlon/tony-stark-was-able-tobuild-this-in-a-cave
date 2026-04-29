package world

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"

	"skyra-v05/src/reality"
)

type Thread struct {
	id     string
	Routes map[RouteKey]*Route
}

type RouteKey struct {
	A, B string
}

func NewRouteKey(a, b string) RouteKey {
	if a < b {
		return RouteKey{A: a, B: b}
	}
	return RouteKey{A: b, B: a}
}

type Route struct {
	Parent  string
	Active  bool
	Entries []reality.Relation
}

func NewThread() *Thread {
	return &Thread{
		id:     "thread",
		Routes: make(map[RouteKey]*Route),
	}
}

func NewThreadID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (t *Thread) ID() string { return t.id }

func (t *Thread) Create(r reality.Relation) reality.Reality {
	return t
}

func (t *Thread) Realize(r reality.Relation) string {
	key := NewRouteKey(r.Origin, r.ID)

	route, ok := t.Routes[key]
	if !ok {
		route = &Route{Parent: r.Origin, Active: true}
		t.Routes[key] = route
	}
	if !route.Active {
		route.Parent = r.Origin
		route.Active = true
	}
	route.Entries = append(route.Entries, r)

	var sb strings.Builder

	ref := extractRef(r.Impulse)
	if ref != "" {
		resolved := t.resolveRef(r.Origin, ref)
		if resolved != "" {
			sb.WriteString(resolved)
		}
	}

	exchange := t.exchangeBetween(r.Origin, r.ID)
	if exchange != "" {
		sb.WriteString("exchange:\n" + exchange)
	}

	returnTo := t.findReturnTarget(r.ID)
	if returnTo != "" {
		sb.WriteString(fmt.Sprintf("return: %s\n", returnTo))
	}

	active := t.activeRoutes(r.ID, r.Origin)
	if active != "" {
		sb.WriteString("active:\n" + active)
	}

	return sb.String()
}

func (t *Thread) CloseRoute(a, b string) {
	key := NewRouteKey(a, b)
	if route, ok := t.Routes[key]; ok {
		route.Active = false
	}
}

func (t *Thread) exchangeBetween(a, b string) string {
	key := NewRouteKey(a, b)
	route, ok := t.Routes[key]
	if !ok {
		return ""
	}
	var sb strings.Builder
	for i, entry := range route.Entries {
		sb.WriteString(fmt.Sprintf("  [%d] %s: %s\n", i, entry.Origin, entry.Impulse))
	}
	return sb.String()
}

func (t *Thread) findReturnTarget(beingID string) string {
	for key, route := range t.Routes {
		if !route.Active {
			continue
		}
		if key.A != beingID && key.B != beingID {
			continue
		}
		if route.Parent == beingID {
			continue
		}
		return route.Parent
	}
	return ""
}

func (t *Thread) activeRoutes(beingID, currentPeer string) string {
	var sb strings.Builder
	for key, route := range t.Routes {
		if !route.Active {
			continue
		}
		if key.A != beingID && key.B != beingID {
			continue
		}
		peer := key.A
		if key.A == beingID {
			peer = key.B
		}
		label := ""
		switch {
		case peer == currentPeer:
			label = "current"
		case route.Parent != beingID:
			label = "waiting on you"
		default:
			label = "you opened"
		}
		sb.WriteString(fmt.Sprintf("  %s — %s, %d entries\n", peer, label, len(route.Entries)))
	}
	return sb.String()
}

func (t *Thread) resolveRef(beingID, ref string) string {
	colon := strings.Index(ref, ":")
	if colon == -1 {
		return ""
	}
	peer := strings.TrimSpace(ref[:colon])
	rng := strings.TrimSpace(ref[colon+1:])

	key := NewRouteKey(beingID, peer)
	route, ok := t.Routes[key]
	if !ok || len(route.Entries) == 0 {
		return ""
	}

	parts := strings.SplitN(rng, "-", 2)
	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return ""
	}
	end := start
	if len(parts) == 2 {
		end, err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return ""
		}
	}
	if start < 0 || start >= len(route.Entries) {
		return ""
	}
	if end >= len(route.Entries) {
		end = len(route.Entries) - 1
	}

	var sb strings.Builder
	sb.WriteString("context (~ref " + ref + "):\n")
	for i, entry := range route.Entries[start : end+1] {
		sb.WriteString(fmt.Sprintf("  [%d] %s: %s\n", i+start, entry.Origin, entry.Impulse))
	}
	return sb.String()
}

func extractRef(impulse string) string {
	idx := strings.Index(impulse, "~ref")
	if idx == -1 {
		return ""
	}
	rest := strings.TrimSpace(impulse[idx+4:])
	end := len(rest)
	for _, delim := range []string{"~", "|"} {
		if i := strings.Index(rest, delim); i != -1 && i < end {
			end = i
		}
	}
	return strings.TrimSpace(rest[:end])
}
