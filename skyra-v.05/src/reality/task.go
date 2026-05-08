package reality

import (
	"fmt"
	"strings"
)

type Task struct {
	Name        string
	Description string
	Assumptions []string
	Commands    []string
	Validation  string
	AcceptedBy  string
	State       string // "active", "held", "done", "dropped"
	Items       []*Task
}

func (t *Task) Realize(indent int, open bool) string {
	prefix := strings.Repeat(" ", indent)
	var sb strings.Builder

	if open {
		sb.WriteString(fmt.Sprintf("%s▾ %s [%s]\n", prefix, t.Name, t.State))
		inner := strings.Repeat(" ", indent+2)
		if t.Description != "" {
			sb.WriteString(inner + t.Description + "\n")
		}
		if len(t.Assumptions) > 0 {
			sb.WriteString(inner + "assumes: " + strings.Join(t.Assumptions, "; ") + "\n")
		}
		if t.Validation != "" {
			sb.WriteString(inner + "done when: " + t.Validation + "\n")
		}
		for _, sub := range t.Items {
			sb.WriteString(sub.Realize(indent+4, false))
		}
	} else {
		sb.WriteString(fmt.Sprintf("%s▸ %s [%s]\n", prefix, t.Name, t.State))
	}

	return sb.String()
}
