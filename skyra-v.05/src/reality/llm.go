package reality

import (
	"fmt"
	"skyra-v05/src/debug"
)

type Provider struct {
	id    string
	Model string
	Call  func(present string) (string, error)
}

func (p *Provider) ID() string { return p.id }

func (p *Provider) Create(r *Relation) Reality {
	return p
}

func (p *Provider) Realize(r *Relation) string {
	present := p.derivePresent(r)
	debug.Log("[provider]: present →")
	debug.Log(present)
	debug.Log("[provider]: calling", p.id)
	response, err := p.Call(present)
	if err != nil {
		fmt.Println("provider error:", err)
		debug.Log("[provider]: error →", err)
		return ""
	}
	debug.Log("[provider]: response →", response)
	return response
}

func (p *Provider) derivePresent(r *Relation) string {
	result := ""
	for _, parser := range r.Parsers {
		result += parser()
	}
	if r.Impulse != "" {
		result += "\nmessage: " + r.Impulse + "\n"
	}
	return result
}

func (p *Provider) Parse() string {
	return ""
}

type LLM struct {
	id        string
	Realities map[string]Reality
}

func NewLLM() *LLM {
	return &LLM{
		id:        "llm",
		Realities: make(map[string]Reality),
	}
}

func (l *LLM) ID() string { return l.id }

func (l *LLM) Create(r *Relation) Reality {
	name, _ := Extract(r.Impulse, "~name", "llm")
	model, _ := Extract(r.Impulse, "~model", "llm")
	l.Realities[name] = &Provider{id: name, Model: model}
	return l
}

func (l *LLM) Realize(r *Relation) string {
	target, ok := l.Realities[r.ID]
	if !ok {
		return ""
	}
	return target.Realize(r)
}

func (l *LLM) Parse() string {
	return ""
}

func (l *LLM) Provider(name string) Reality {
	return l.Realities[name]
}

func (l *LLM) WireCall(name string, call func(string) (string, error)) {
	p, ok := l.Realities[name]
	if !ok {
		return
	}
	p.(*Provider).Call = call
}
