package reality

import (
	"fmt"
	"skyra-v05/src/debug"
	"strings"
)

type Context struct {
	id     string
	Owner  string
	Memory *Memory
	LLM    Reality
	Warm   map[string][]*MemNode
	Active string
}

func (c *Context) ID() string { return c.id }

func (c *Context) Create(r *Relation) Reality {
	return &Context{id: "context", Warm: make(map[string][]*MemNode)}
}

func (c *Context) Realize(r *Relation) string {
	return ""
}

func (c *Context) Heat(relationship string) {
	log := func(args ...any) { debug.Being(c.Owner, "context", args...) }

	if c.Active == relationship && len(c.Warm[relationship]) > 0 {
		return
	}

	c.Active = relationship

	var loaded []*MemNode
	seen := map[string]bool{}

	for _, entity := range c.Memory.Graph.Entities {
		memories := c.Memory.Graph.MemoriesForEntity(entity.ID)
		for _, node := range memories {
			if node.Relationship != relationship || seen[node.ID] {
				continue
			}
			seen[node.ID] = true
			loaded = append(loaded, node)
		}
	}

	c.Warm[relationship] = loaded
	log("[context]: heated", relationship, "—", len(loaded), "memories")
}

func (c *Context) Parse(relationship string) string {
	nodes := c.Warm[relationship]
	if len(nodes) == 0 {
		return ""
	}

	understandings := filterByType(nodes, "understanding")
	tensions := filterByType(nodes, "tension")

	if len(understandings) == 0 && len(tensions) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("memory context:\n")
	for _, n := range understandings {
		sb.WriteString("  [understanding] " + n.Content + "\n")
	}
	for _, n := range tensions {
		sb.WriteString("  [tension] " + n.Content + "\n")
	}
	return sb.String()
}

func filterByType(nodes []*MemNode, artifactType string) []*MemNode {
	var result []*MemNode
	for _, n := range nodes {
		if n.Type == artifactType {
			result = append(result, n)
		}
	}
	return result
}

func (c *Context) Store(content, relationship, artifactType string) string {
	log := func(args ...any) { debug.Being(c.Owner, "context", args...) }

	if c.LLM == nil {
		log("[context]: no llm, storing raw")
		msg := c.Memory.Store(content, relationship, artifactType, EdgeLayer{Type: "episode", Weight: 1.0})
		c.invalidate(relationship)
		return msg
	}

	existing := c.Memory.Query(content, relationship, "")

	lr := &Relation{
		Impulse: content,
		Parsers: make(map[string]Parser),
	}

	lr.Attach("system", c.system)
	lr.Attach("context-job", func() string {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("artifact type: %s\n", artifactType))
		sb.WriteString(fmt.Sprintf("relationship: %s\n", relationship))
		sb.WriteString(fmt.Sprintf("content to store:\n%s\n", content))
		if existing != "no relevant memories found" {
			sb.WriteString(fmt.Sprintf("\nexisting memories in this relationship:\n%s\n", existing))
		}
		return sb.String()
	})

	result := c.LLM.Realize(lr)
	log("[context]: curated →", result)

	cleaned := content
	if cl, err := ExtractTag(result, "content"); err == nil {
		cleaned = cl
	}

	finalType := artifactType
	if t, err := ExtractTag(result, "type"); err == nil {
		if isValidArtifactType(t) {
			finalType = t
		}
	}

	var entities []string
	if e, err := ExtractTag(result, "entities"); err == nil {
		for _, entity := range strings.Split(e, ",") {
			entity = strings.TrimSpace(entity)
			if entity != "" {
				entities = append(entities, entity)
				c.Memory.Extractor.Learn(entity)
			}
		}
	}

	action := "store"
	if a, err := ExtractTag(result, "action"); err == nil {
		action = strings.TrimSpace(a)
	}

	judgment := "complement"
	if j, err := ExtractTag(result, "judgment"); err == nil {
		judgment = strings.TrimSpace(j)
	}

	switch action {
	case "discard":
		log("[context]: discarded as redundant")
		return "already known"
	}

	switch judgment {
	case "supersede":
		if target, err := ExtractTag(result, "supersedes"); err == nil {
			c.decayByContent(target, relationship)
			log("[context]: superseded existing memory")
		}
	case "contradict":
		if target, err := ExtractTag(result, "contradicts"); err == nil {
			c.markTension(target, relationship)
			log("[context]: marked contradiction as tension")
		}
		finalType = "tension"
	}

	msg := c.Memory.Store(cleaned, relationship, finalType, EdgeLayer{Type: "episode", Weight: 1.0})
	c.invalidate(relationship)

	if len(entities) > 0 {
		log("[context]: learned entities:", entities)
	}

	return msg
}

func (c *Context) decayByContent(content, relationship string) {
	for _, node := range c.Memory.Graph.Nodes {
		if node.Relationship == relationship && strings.Contains(node.Content, content) {
			node.Weight *= 0.3
		}
	}
}

func (c *Context) markTension(content, relationship string) {
	for _, node := range c.Memory.Graph.Nodes {
		if node.Relationship == relationship && strings.Contains(node.Content, content) {
			node.Type = "tension"
		}
	}
}

func (c *Context) invalidate(relationship string) {
	delete(c.Warm, relationship)
}

func (c *Context) Retrieve(query, relationship, artifactType string) string {
	c.Heat(relationship)
	return c.Memory.Query(query, relationship, artifactType)
}

func (c *Context) system() string {
	return `You are a memory curator. Your job is to evaluate new memories against existing ones before storage.

You receive content that a being wants to remember, along with existing memories in the same relationship.

Your responsibilities:
1. Clean the content — tighten language, strip noise, keep signal
2. Extract entities — names, concepts, tools, patterns mentioned in the content
3. Judge the relationship to existing memories:
   - supersede: the new memory replaces an old understanding. Name what it supersedes.
   - complement: the new memory adds nuance without replacing. Both keep their weight.
   - contradict: the new memory conflicts with an existing one. The old becomes a tension.
   - discard: this is already fully covered by existing memories.
4. Classify — confirm or reclassify the artifact type (trace, salience, tension, understanding)

Respond with exactly these tags:

<content>the cleaned content to store</content>
<type>trace|salience|tension|understanding</type>
<entities>comma,separated,entity,names</entities>
<judgment>supersede|complement|contradict</judgment>
<action>store|discard</action>

If judgment is supersede, also include:
<supersedes>the content of the memory being replaced</supersedes>

If judgment is contradict, also include:
<contradicts>the content of the conflicting memory</contradicts>

If discarding, still include all tags but set content to empty.
Keep your response to just the tags. No commentary.`
}

func isValidArtifactType(t string) bool {
	switch t {
	case "trace", "salience", "tension", "understanding":
		return true
	}
	return false
}

func findContext(r *Relation) *Context {
	if r.Realities == nil {
		return nil
	}
	if c, ok := r.Realities["context"]; ok {
		if ctx, ok := c.(*Context); ok {
			return ctx
		}
	}
	return nil
}
