package kernel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PromptRegistry struct {
	prompts map[string]string
}

func LoadPromptRegistry(root string) (*PromptRegistry, error) {
	files := map[string]string{
		"experience": filepath.Join(root, "experience", "experience.txt"),
		"understand": filepath.Join(root, "understand", "understandmenu.txt"),
		"reference":  filepath.Join(root, "understand", "interpret", "reference", "Reference.text"),
		"infer":      filepath.Join(root, "understand", "interpret", "infer", "infer.txt"),
		"resolve":    filepath.Join(root, "understand", "interpret", "resolve", "resolve.txt"),
		"interact":   filepath.Join(root, "interact", "interact.txt"),
	}

	registry := &PromptRegistry{prompts: make(map[string]string, len(files))}
	for name, path := range files {
		body, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("load prompt %q: %w", name, err)
		}
		registry.prompts[name] = string(body)
	}
	return registry, nil
}

func (r *PromptRegistry) Render(name string, values map[string]string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("prompt registry is nil")
	}
	raw, ok := r.prompts[name]
	if !ok {
		return "", fmt.Errorf("prompt %q not found", name)
	}

	out := raw
	for key, value := range values {
		out = strings.ReplaceAll(out, "{"+key+"}", value)
		out = strings.ReplaceAll(out, "<"+key+">", value)
	}
	return out, nil
}
