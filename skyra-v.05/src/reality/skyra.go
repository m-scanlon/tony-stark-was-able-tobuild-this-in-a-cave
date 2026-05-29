package reality

// Skyra DNA — the first expression of the Reality interface.
// Every field here is the same as every other DNA. The behavior
// comes from what the genome registers in ObserveFns/ExpressFns
// and where this Reality sits in the topology.
type Skyra struct {
	id             string  // unique per instance, assigned at creation
	Weight         float64 // global weight — EMA across all traversals
	TraversalCount int     // proper time — how many relations have passed through
	LastTraversed  int     // traversal count when last activated
	Alpha          float64 // EMA smoothing factor — slow = stable, fast = reactive
	Type           string  // current expression — set by Competence each frame
	Content        string  // what this cell holds — identity for PFC, memory for memory, skill for skill
	Description    string  // what this cell type is responsible for — shapes traversal, not shown during inference

	// the topology — where this Reality sits in the graph
	Relationships map[string]Reality // descent targets — what this node knows
	Expressors    map[string]Reality // ascent targets — what this node can do

	// the DNA — keyed by Type, selected by Competence
	ObserveFns map[string]func(*Relation)            // intake behaviors per type — needed: pass-through, pick-up, propagate
	ExpressFns map[string]func(*Relation) *Relation  // output behaviors per type
	// expressors needed: motor (act), thought (provider call), memory (store/retrieve), specialist (promoted dense region)
}

func (s *Skyra) ID() string {
	return s.id
}

// Create — instantiate a new cell with this DNA. The genome calls this.
func (s *Skyra) Create(r *Relation) Reality {
	return &Skyra{
		id:            NewID(),
		Weight:        1.0,
		Alpha:         0.1,
		Relationships: make(map[string]Reality),
		Expressors:    make(map[string]Reality),
		ObserveFns:    make(map[string]func(*Relation)),
		ExpressFns:    make(map[string]func(*Relation) *Relation),
	}
}

// Competence — the cell reads itself before processing.
// Computes activation from its own weights and the Relation's
// binding fields, then sets Type. This gates which ObserveFn
// and ExpressFn fire this frame.
func (s *Skyra) Competence(r *Relation) {
	activation := s.Activation(r)

	// thresholds, type names, and number of types are all placeholders
	switch {
	case activation > 0.8:
		s.Type = "integrator"
	case activation > 0.6:
		s.Type = "specialist" // episodic processor: batch extraction, memory deposit
	case activation > 0.4:
		s.Type = "processor"
	case activation > 0.2:
		s.Type = "context"
	default:
		s.Type = "dormant"
	}
}

// Realize — one traversal through this cell.
// Competence → Observe → descent → ascent → Express → weight update.
// The Relation enters and leaves transformed. So does the cell.
func (s *Skyra) Realize(r *Relation) *Relation {
	if r.Cancelled {
		s.Decay()
		return r
	}
	if r.Visited[s.id] {
		return r
	}
	r.mu.Lock()
	r.Visited[s.id] = true
	r.mu.Unlock()

	// first competence read — what am I on intake?
	s.Competence(r)

	// descent phase — the cell absorbs from the signal
	s.Observe(r)

	// deeper into the topology — promoted nodes first, then memories
	for _, rel := range s.Relationships {
		r = rel.Realize(r)
	}

	// decay unvisited neighbors — cancelled traversal, one level deep
	for id, rel := range s.Relationships {
		if !r.Visited[id] {
			rel.Realize(&Relation{Cancelled: true})
		}
	}

	// ascent phase — expressors fire on the way back up
	for _, e := range s.Expressors {
		r = e.Realize(r)
	}

	// second competence read — what am I on output, after descent changed things?
	s.Competence(r)

	// the cell outputs
	r = s.Express(r)

	// the cell ages — proper time increments, weights update
	s.TraversalCount++
	s.LastTraversed = s.TraversalCount
	s.Reinforce()

	return r
}

// Observe — mutual observation. The Relation observes itself through
// the encounter first, then the cell's type-specific intake fires.
func (s *Skyra) Observe(r *Relation) {
	r.mu.Lock()
	r.LastSeen = s
	r.mu.Unlock()
	r.Observe(r)
	if s.Content != "" && s.Activation(r) > 0.2 {
		r.mu.Lock()
		r.Thoughts = append(r.Thoughts, s.Content)
		r.Load += len(s.Content)
		r.Signal -= 0.01 // TODO: cost should be based on token count of content, not flat rate
		r.mu.Unlock()
	}
	if fn, ok := s.ObserveFns[s.Type]; ok {
		fn(r)
	}
}

// Express — mutual expression. The Relation expresses itself first,
// then the cell's type-specific output fires. May call a Provider
// from r.Providers if the cell needs inference.
func (s *Skyra) Express(r *Relation) *Relation {
	r = r.Express(r)
	if fn, ok := s.ExpressFns[s.Type]; ok {
		r = fn(r)
	}
	return r
}

// Activation — base weight * recency * product of all binding fields.
// This is the score that Competence reads to determine what the cell is.
func (s *Skyra) Activation(rel *Relation) float64 {
	base := s.Weight * Recency(s.TraversalCount, s.LastTraversed)
	if field, ok := rel.Fields[s.id]; ok {
		for _, w := range field {
			base *= w
		}
	}
	return base
}

// Reinforce — EMA update on the return path. The cell was traversed,
// so its weight moves toward 1.0. Alpha controls how fast.
func (s *Skyra) Reinforce() {
	s.Weight = s.Alpha*1.0 + (1-s.Alpha)*s.Weight
}

// Decay — EMA update when NOT traversed. Weight drifts toward 0.
// Starvation is built into the formula.
func (s *Skyra) Decay() {
	s.Weight = (1 - s.Alpha) * s.Weight
}
