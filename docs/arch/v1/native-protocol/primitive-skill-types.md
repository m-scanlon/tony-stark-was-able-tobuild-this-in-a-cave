# Primitive Skill Types (Native Protocol v1)

All execution is a skill instantiated as a job. This registry lists native protocol primitives only.

---

## Primitive Registry

| Primitive | Trigger | Live user? |
|---|---|---|
| `user_event` | Ingress receives user input | Yes |
| `reply` | Skill emits user-visible output | Yes |
| `remember` | Workflow commits memory artifacts | No |
| `retrieve` | Context retrieval requested | No |
| `decide` | Candidate-set decision step requested | No |
| `recurse` | Child workflow call requested | No |

---

## Descriptions

### `user_event`
Input primitive interface for user-originated events.

- Command shape: `octos user_event [args...]`
- Implementations: `--impl=voice | text` (kernel allowlist enforced)
- Purpose: normalize inbound user signal into a bounded, typed event for skill execution

### `reply`
Output primitive for user-visible updates and finals.

- Command shape: `octos reply [args...]`
- Used for `update | final | error` response emission

See: `docs/arch/v1/native-protocol/primitives/reply.md`.

### `remember`
Memory-write primitive interface for bounded refs and namespace-targeted persistence.

- Command shape: `octos remember [args...]`
- Supports implementation selection (for example `episodic | semantic | procedural`)
- Enforces layer policy (`working | committed`)

See: `docs/arch/v1/native-protocol/primitives/remember.md`.

### `retrieve`
Bounded context retrieval primitive.

- Command shape: `octos retrieve [args...]`
- Supports implementation selection (for example `ctor | frame | graph`)

See: `docs/arch/v1/native-protocol/primitives/retrieve/retrieve.md`.

### `decide`
Bounded decision primitive over candidate refs.

- Command shape: `octos decide [args...]`
- Returns explicit action (`keep | drop | split | delegate | close`)

See: `docs/arch/v1/native-protocol/primitives/decide.md`.

### `recurse`
Bounded child-call primitive with deterministic return-to-parent contract.

- Command shape: `octos recurse [args...]`
- Enforces depth/context/return budgets and overflow policy

See: `docs/arch/v1/native-protocol/primitives/recurse.md`.

---

## Not Primitive Types

The following are system capabilities/workflows, not native protocol primitives:

- `conversation`
- `synthesize`
- `internet_search`
- `shard_registration`
- `cron`
- `batch`
- `skill_acquisition`
- `compaction`
