# Skill: Integrate

The integrate skill connects the mini graph produced by the reasoning skill to the existing graph. It runs after reasoning completes.

**Triggered by**: reasoning skill completion
**Layer**: observational — Skyra writes freely, no user gate
**Type**: system primitive skill — pre-provisioned in Redis at boot

---

## The Job

For each new observational actor, find where it fits in the existing graph. Most of the time this is strengthening existing structure — updating weights, adding missing edges — not creating new actors.

The graph is large. It never loads into context. Integrate works systematically, one actor at a time, using search to constrain the candidate set.

---

## The Two Signals

Integration confidence is derived from two signals combined:

**1. Semantic similarity** — vector similarity between the new actor and existing actor content. Finds candidates with similar meaning.

**2. Structural similarity** — shared neighbors in the graph. If a candidate actor is already connected to the same entities as the new actor, that's strong evidence they're related.

| Semantic | Structural | Interpretation |
|---|---|---|
| High | High | Strong alias candidate — likely the same entity |
| High | Low | Related but probably distinct — write a `relates_to` edge |
| Low | High | Worth examining — could be coincidence |
| Low | Low | No connection — leave it |

---

## The Process

```
for each new observational actor:
  → skyra search -scope graph -query actor.content
  → returns top N semantically similar existing actors
  → check structural similarity: shared neighbors
  → reason over candidates
      → same entity?  → alias_of edge, update weights, add missing edges
      → related?      → relates_to edge
      → unrelated?    → no action
```

The existing graph never loads into context. Only the candidate set does.

---

## Outputs

Integrate does not create actors. The primary outputs are:

- **`alias_of` edges** — new actor and existing actor are the same entity. Incoming actor is not deleted — the edge records the resolution.
- **New edges** — relationships implied by the new actor that don't exist yet in the graph
- **Weight updates** — existing edges strengthened by new corroborating signal

```
// Mike (existing committed) and Mike (incoming observational) — same entity
skyra write_edge \
  -from "entity:mike_incoming" -to "entity:mike" \
  -type alias_of -weight 1.0 \
  -reasoning "skyra write_edge ... -m \"high semantic + structural similarity. incoming 'Mike' resolves to existing committed entity.\""

// New edge implied by incoming actor that doesn't exist yet
skyra write_edge \
  -from "entity:mike" -to "entity:late_night" \
  -type works_at -weight 0.7 \
  -reasoning "skyra write_edge ... -m \"edge does not exist in graph. adding from incoming signal.\""
```

---

## Rules

- Skyra cannot mutate committed actors. She can add edges to them freely.
- Incoming observational actors are never deleted — alias_of edges record the resolution.
- Every edge written gets a `reasoning` field.
- If confidence is low, do nothing. A missing edge is better than a wrong one.
- We are not designing around correctness. Observational errors don't corrupt the committed layer.

---

## State Contract

Writes to the observational layer only. No user gate. No committed writes.

---

## Validation Criteria

- Every new actor has been evaluated — no actor left unprocessed
- Every edge written has a `reasoning` field
- No committed actors mutated

---

## Skill Contract

```
skill: integrate
tasks:
  1. for each new actor: search + reason + write
boundary_rules:
  search:                       allow_always
  write_edge:                   allow_always (observational layer + edges to committed actors)
  write_actor:                   deny (integrate does not create actors)
  write_actor (skill type):      deny
state_contract: working (observational only, no approval required)
severity_policy:
  low confidence candidate: do nothing — skip
  missing neighbor data:    log and continue
replan_budget: 2
```

---

## Design Note

This skill is designed around the output of the reasoning skill. The quality of integration depends on the quality of the mini graph coming in — atomic actor labels, well-formed edges, accurate reasoning fields. Build and test reasoning first. Revisit integrate when real graph output exists.

---

## Related

- `docs/arch/v1/skill/skill-reasoning.md` — produces the mini graph integrate consumes
- `docs/arch/v1/memory/memory-structure.md` — actor + edge schema, two-tier graph
- `docs/arch/v1/skill/skill.md` — skill schema, closed for modification
- `docs/arch/v1/kernel/kernel.md` — job execution, primitive skills
