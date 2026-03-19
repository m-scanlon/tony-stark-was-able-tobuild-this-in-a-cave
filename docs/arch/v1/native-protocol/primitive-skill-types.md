# Secondary Primitive Skill Types (Native Protocol v1)

All execution is a skill instantiated as a job.

The top-level process model is defined in `thought-processes.md`.

The top-level execution primitive set is:

- `stimuli`
- `process`
- `interact`

A primitive must participate in execution.
These top-level primitives do.
`process` can contain child forms such as `thought`.
`thought` can contain child forms such as `repair`.
`interact` can contain child forms such as `agree` and `reply`.
`agree` may be carried by `reply`.

This registry lists the second-tier native protocol primitives that sit under
that model.

---

## Secondary Primitive Registry

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

- Command shape: `skyra user_event [args...]`
- Implementations: `--impl=voice | text` (kernel allowlist enforced)
- Purpose: normalize inbound user signal into a bounded, typed event for skill execution

### `reply`
Output primitive for user-visible updates and finals.

- Command shape: `skyra reply [args...]`
- Used for `update | final | error` response emission

See: `docs/arch/v1/native-protocol/primitives/reply.md`.

### `remember`
Memory-write primitive interface for bounded refs and namespace-targeted persistence.

- Command shape: `skyra remember [args...]`
- Supports implementation selection (for example `synthesize_append | append_only`)
- Enforces layer policy (`working | committed`)
- Default operational flow is `synthesize -> append`.

See: `docs/arch/v1/native-protocol/primitives/remember.md`.

### `retrieve`
Bounded context retrieval primitive.

- Command shape: `skyra retrieve [args...]`
- Supports implementation selection (for example `ctor | frame | graph`)

See: `docs/arch/v1/native-protocol/primitives/retrieve/retrieve.md`.

### `decide`
Bounded decision primitive over candidate refs.

- Command shape: `skyra decide [args...]`
- Returns explicit action (`keep | drop | split | delegate | close`)

See: `docs/arch/v1/native-protocol/primitives/decide.md`.

### `recurse`
Bounded child-call primitive with deterministic return-to-parent contract.

- Command shape: `skyra recurse [args...]`
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
