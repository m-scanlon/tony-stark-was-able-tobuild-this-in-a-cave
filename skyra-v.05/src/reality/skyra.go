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

	// the topology — where this Reality sits in the graph
	Relationships map[string]Reality // descent targets — what this node knows
	Expressors    map[string]Reality // ascent targets — what this node can do

	// the DNA — keyed by Type, selected by Competence
	ObserveFns map[string]func(*Relation)            // intake behaviors per type
	ExpressFns map[string]func(*Relation) *Relation  // output behaviors per type
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
	case activation > 0.5:
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
	if r.Visited[s.id] {
		return r
	}
	r.Visited[s.id] = true

	// the cell reads itself — what am I right now?
	s.Competence(r)

	// descent phase — the cell absorbs from the signal
	s.Observe(r)

	// deeper into the topology — promoted nodes first, then memories
	for _, rel := range s.Relationships {
		r = rel.Realize(r)
	}

	// ascent phase — expressors fire on the way back up
	for _, e := range s.Expressors {
		r = e.Realize(r)
	}

	// the cell outputs — may call a Provider from the Relation
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
	r.Observe(r)
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
	for _, field := range rel.Fields {
		if w, ok := field[s.id]; ok {
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
