# Native Protocol (v1)

This protocol is the object-oriented representation of system execution.
The kernel is the interpreter.

This is the native command protocol for skill execution and skill composition.

## Object Model

1. Skill = class
2. Job = object instance of a skill
3. Nested skill = callable composition target declared inside a parent skill
4. Tasks do not exist in this protocol model

Everything is composed by skills calling skills.

## Base Schema

```
octos <skill> [args...]
```

One command. One root skill invocation.

## Constructor (Context Nodes)

A skill constructor is the set of context nodes used by the skill when solving a problem.

Constructor inputs are passed as node references:

```
--ctor.<param>=node:<node_id>
```

Examples:

```
octos planner --ctor.user=node:user.mike --ctor.goal=node:goal.fitness --query="best next action"
```

## Composition Forms

Nested skill call:

```
octos <root_skill>.<nested_skill> [args...]
```

Skill-as-input:

```
--skill.<param>=skill:<skill_id>
```

Nested-skill provisioning declaration (inside skill contract):

```
provisioned_skills: [<nested_skill_1>, <nested_skill_2>, ...]
```

Reasoning scope rule:

- While executing a root skill, the model may reason over only the nested skills listed in that root skill's `provisioned_skills[]`.

## Definition Requirements

A composed call is valid only if all are true:

1. The root skill is defined and committed.
2. The nested skill is defined and committed.
3. Both skills are provisioned in Redis (executable).
4. The nested skill is listed in the root skill's `provisioned_skills[]`.
5. Boundary policy allows root -> nested invocation (allowed or gated, not denied).
6. The call remains inside the originating user intent scope. No exceptions.

If any condition fails, gateway/kernel rejects the call.

Composition examples:

```
octos orchestrator.search --query="today's booking"
octos integrate --skill.source=skill:reasoning.v1 --skill.target=skill:search.v3
```

## Decision Continuity and Job Closure

The protocol assumes models can forget. Decision continuity is enforced by system state, not model memory.

Decision record (append-only):

```
decision_record {
  decision_id
  job_id
  call_id
  must_call[]      // required skill calls before close
  must_not_call[]  // forbidden skill calls
  close_allowed    // true | false
  reason
  created_at
}
```

Rules:

1. Every material decision appends a `decision_record`.
2. Downstream skill calls are validated against active `must_not_call[]`.
3. Job closure is explicit, never implicit:

```
octos job.close --job=<job_id> --decision=<decision_id> --status=done|abandoned|superseded
```

4. Close is rejected if unresolved `must_call[]` obligations remain (unless explicitly waived by policy/user approval).

## Notes

- `skill:<skill_id>` is the canonical skill reference format.
- `node:<node_id>` is the canonical constructor node reference format.
- "method" may be used informally to describe the relationship from root skill to nested skill, but it is not a protocol primitive.
- Command syntax stays native CLI; shards do not emit JSON envelopes.
- Detailed lineage/idempotency/depth/error contract lives in:
  `docs/arch/v1/skill-composition-protocol.md`.
