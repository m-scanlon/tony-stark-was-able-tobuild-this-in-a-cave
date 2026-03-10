# Skill: Reasoning

The reasoning skill is how Skyra builds her memory. It runs as a background job when the user is offline. It takes raw session data and turns it into graph nodes and edges in the observational layer.

**Triggered by**: Cron Service
**Layer**: observational — Skyra writes freely, no user gate
**Type**: system primitive skill — pre-provisioned in Redis at boot

---

## The Job

Two sequential tasks. Task 1 must complete before Task 2 runs.

```
cron fires
  → reasoning job created
  → task 1: decompose
  → task 2: relate
```

---

## Task 1 — Decompose

**Input**: unprocessed session history + VAD time series, married on `turn_id`

**Goal**: extract every entity, theme, and atomic fact from the session. One node per thing. No sentences. No relationships yet — that's task 2.

**Output**: observational nodes written to the graph

```
skyra write_node -type entity -layer observational \
  -content "late night" \
  -reasoning "skyra write_node ... -m \"repeated across 6 sessions, user active after 11pm\""

skyra write_node -type entity -layer observational \
  -content "Skyra project" \
  -reasoning "skyra write_node ... -m \"primary topic of session, referenced 14 times\""
```

**Rules**:
- One node per entity. Do not merge entities — keep them atomic. When in doubt, create separate nodes. Integrate will resolve duplicates via `alias_of` edges. False positives are better than merging distinct entities.
- `content` is the label for the thing. Not a sentence. Not a relationship.
- Every node gets a `reasoning` field. Skyra explains why she created it.
- Confidence score derived from frequency + affect signal in the session.

---

## Task 2 — Relate

**Input**: nodes produced by task 1 + existing graph (committed + observational)

**Goal**: reason over relationships between the new nodes and the existing graph. Write edges.

**Output**: observational edges written to the graph

```
skyra write_edge \
  -from "entity:mike" -to "entity:skyra_project" \
  -type works_on -weight 0.9 \
  -reasoning "skyra write_edge ... -m \"14 references to skyra project across session, high arousal, positive valence\""

skyra write_edge \
  -from "entity:mike" -to "entity:late_night" \
  -type works_at -weight 0.7 \
  -reasoning "skyra write_edge ... -m \"session timestamps consistently 11pm-2am, VAD arousal elevated\""

skyra write_edge \
  -from "entity:skyra_project" -to "entity:late_night" \
  -type associated_with -weight 0.6 \
  -reasoning "skyra write_edge ... -m \"co-occurrence — skyra project work clusters with late night sessions\""
```

**Rules**:
- Check for existing nodes before creating relationships to duplicates. Use `alias_of` if an existing node represents the same entity.
- Skyra can add edges to committed nodes freely.
- Every edge gets a `reasoning` field.
- Weight derived from frequency, affect magnitude, and co-occurrence signal.
- Edge types are not exhaustive — new types can emerge.

---

## State Contract

Writes to the observational layer only. No user gate. No committed writes.

---

## Validation Criteria

- Every entity mentioned in session history has a corresponding node
- Every node has a `reasoning` field
- Every edge has a `reasoning` field
- No node `content` is a sentence — atomic labels only
- No duplicate nodes for the same entity

---

## Skill Contract

```
skill: reasoning
tasks:
  1. decompose  — session history + VAD → observational nodes
  2. relate     — nodes + existing graph → observational edges
boundary_rules:
  write_node (non-skill types): allow_always (observational layer)
  write_node (skill type):      deny — use update_skill
  write_edge:                   allow_always (observational layer + edges to committed nodes)
  read_graph:                   allow_always
state_contract: working (observational only, no approval required)
severity_policy:
  duplicate entity detected: adjust locally — resolve to existing node
  missing session data:      log and continue — partial pass is valid
replan_budget: 2
```

---

## Related

- `docs/arch/v1/memory-structure.md` — node + edge schema, two-tier graph, cron pass
- `docs/arch/v1/observational-store.md` — four streams, VAD, session history
- `docs/arch/v1/skill-lifecycle.md` — how skills crystallize from observational nodes
- `docs/arch/v1/kernel.md` — cron service, job execution
