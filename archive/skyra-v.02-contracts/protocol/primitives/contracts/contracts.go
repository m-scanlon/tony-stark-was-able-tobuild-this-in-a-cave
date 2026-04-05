package primitives

import "time"

type PrimitiveName string

const (
	PrimitiveRecall  PrimitiveName = "recall"
	PrimitiveLearn   PrimitiveName = "learn"
	PrimitiveObserve PrimitiveName = "observe"
	PrimitiveAct     PrimitiveName = "act"
)

type ObserveArgs struct {
	Target     string         `json:"target,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

type ObserveResult struct {
	Data map[string]any `json:"data,omitempty"`
}

type RetainedArtifactKind string

const (
	RetainedArtifactTrace         RetainedArtifactKind = "trace"
	RetainedArtifactUnderstanding RetainedArtifactKind = "understanding"
	RetainedArtifactSalience      RetainedArtifactKind = "salience"
	RetainedArtifactTension       RetainedArtifactKind = "tension"
)

// PrimitiveInvocation is an internal parsed execution record produced after the
// kernel validates a minimal protocol.StimulusEnvelope and resolves caller/target.
type PrimitiveInvocation struct {
	CommandID    string        `json:"command_id"`
	CallingActor string        `json:"calling_actor,omitempty"`
	TargetActor  string        `json:"target_actor"`
	EpisodeID    string        `json:"episode_id,omitempty"`
	IntentID     string        `json:"intent_id,omitempty"`
	Primitive    PrimitiveName `json:"primitive"`
	Reason       string        `json:"reason"`
	EmittedAt    time.Time     `json:"emitted_at"`
}

type PrimitiveResultEvent struct {
	CommandID    string        `json:"command_id"`
	CallingActor string        `json:"calling_actor,omitempty"`
	TargetActor  string        `json:"target_actor"`
	EpisodeID    string        `json:"episode_id,omitempty"`
	IntentID     string        `json:"intent_id,omitempty"`
	Primitive    PrimitiveName `json:"primitive"`
	ResultKind   string        `json:"result_kind"`
	CompletedAt  time.Time     `json:"completed_at"`
}

type RecallQueryKind string

const (
	RecallQueryEntity       RecallQueryKind = "entity"
	RecallQueryRelationship RecallQueryKind = "relationship"
	RecallQueryBundle       RecallQueryKind = "bundle"
)

type RecallQuery struct {
	Kind           RecallQueryKind `json:"kind"`
	EntityID       string          `json:"entity_id,omitempty"`
	RelationshipID string          `json:"relationship_id,omitempty"`
	LeftEntityID   string          `json:"left_entity_id,omitempty"`
	RightEntityID  string          `json:"right_entity_id,omitempty"`
}

type RecallArgs struct {
	Queries []RecallQuery `json:"queries"`
	TopK    int           `json:"top_k,omitempty"`
}

type RecalledArtifact struct {
	ArtifactID             string               `json:"artifact_id"`
	Kind                   RetainedArtifactKind `json:"kind"`
	Score                  float64              `json:"score"`
	MatchedEntityIDs       []string             `json:"matched_entity_ids,omitempty"`
	MatchedRelationshipIDs []string             `json:"matched_relationship_ids,omitempty"`
}

type RecallPackage struct {
	RetainedArtifactIDs []string           `json:"retained_artifact_ids"`
	Matches             []RecalledArtifact `json:"matches"`
}

type RecallInvocation struct {
	PrimitiveInvocation
	Args RecallArgs `json:"args"`
}

type RecallResultEvent struct {
	PrimitiveResultEvent
	Result RecallPackage `json:"result"`
}

type LearnArgs struct {
	EpisodeID string `json:"episode_id"`
}

type LearnPackage struct {
	EpisodeID           string   `json:"episode_id"`
	RetainedArtifactIDs []string `json:"retained_artifact_ids"`
	StructureUpdateIDs  []string `json:"structure_update_ids,omitempty"`
}

type LearnInvocation struct {
	PrimitiveInvocation
	Args LearnArgs `json:"args"`
}

type LearnResultEvent struct {
	PrimitiveResultEvent
	Result LearnPackage `json:"result"`
}

type ObserveInvocation struct {
	PrimitiveInvocation
	Args ObserveArgs `json:"args"`
}

type ObserveResultEvent struct {
	PrimitiveResultEvent
	Result ObserveResult `json:"result"`
}

type ActArgs struct {
	Target    string `json:"target"`
	Content   string `json:"content"`
	Modality  string `json:"modality"`
	Timestamp string `json:"timestamp"`
}

type ActResult struct {
	Data map[string]any `json:"data,omitempty"`
}

type ActInvocation struct {
	PrimitiveInvocation
	Args ActArgs `json:"args"`
}

type ActResultEvent struct {
	PrimitiveResultEvent
	Result ActResult `json:"result"`
}
