package reality

import (
	"fmt"
	"strings"
)

type Economics struct {
	id     string
	Fields map[string]int
}

func NewEconomics() *Economics {
	return &Economics{
		id:     "economics",
		Fields: make(map[string]int),
	}
}

func (e *Economics) ID() string { return e.id }

func (e *Economics) Create(r *Relation) Reality {
	return e
}

func (e *Economics) Set(field string, value int) {
	e.Fields[field] = value
}

func (e *Economics) Realize(r *Relation) string {
	if len(e.Fields) == 0 {
		return ""
	}
	var sb strings.Builder
	for field, value := range e.Fields {
		sb.WriteString(fmt.Sprintf("%s: %d\n", field, value))
	}
	return sb.String()
}
