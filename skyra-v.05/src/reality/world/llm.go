package world

import (
	"skyra-v05/src/reality"
	"skyra-v05/src/reality/being"
)

type Provider struct {
	id    string
	Model string
	Call  func(present string) (string, error)
}

func (p *Provider) ID() string { return p.id }

func (p *Provider) Create(r reality.Relation) reality.Reality {
	return p
}

func (p *Provider) Realize(r reality.Relation) string {
	response, err := p.Call(r.Impulse)
	if err != nil {
		return ""
	}
	return response
}

type LLM struct {
	World
}

func NewLLM() *LLM {
	return &LLM{
		World: World{
			id:        "llm",
			name:      "llm",
			Realities: make(map[string]reality.Reality),
		},
	}
}

func (l *LLM) Register(name, model string, call func(string) (string, error)) {
	l.Realities[name] = &Provider{id: name, Model: model, Call: call}
}

func (l *LLM) Create(r reality.Relation) reality.Reality {
	name, _ := being.Extract(r.Impulse, "~name", "llm")
	model, _ := being.Extract(r.Impulse, "~model", "llm")
	l.Realities[name] = &Provider{id: name, Model: model}
	return l
}

func (l *LLM) Device(name string) Device {
	p, ok := l.Realities[name]
	if !ok {
		return nil
	}
	return p.(Device)
}

func (l *LLM) WireCall(name string, call func(string) (string, error)) {
	p, ok := l.Realities[name]
	if !ok {
		return
	}
	p.(*Provider).Call = call
}
