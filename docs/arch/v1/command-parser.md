# Command Parser

## Syntax

```
skyra <agent> <tool> [args]
```

`me` is a reserved agent name — execute on the current shard, no heap re-entry.

```
skyra me <tool> [args]
```

## Examples

```
skyra me play_audio --file alert.wav
skyra samsung_qled_65 turn_on
skyra music_agent play --query "chill playlist"
skyra gym_agent log_workout --type run --duration 30
skyra mac_mini run_job --input "reorganize filesystem"
```

## Resolution

Skyra always resolves a skill call by walking:

```
1. Where am I?                → location from ingress shard fingerprint
2. What agents are here?      → query agent skill registry filtered by location
3. What does the user mean?   → match intent to an agent skill
4. Which agent holds it?      → dispatch
```

If the user names an agent directly, steps 1–3 still apply — the resolver confirms the agent is reachable and holds the requested skill before dispatching.

## Routing

```
skyra me <tool>       → local execution, no heap re-entry
skyra <agent> <tool>  → heap re-entry → External Router → agent's shard
```

## Related

- `docs/arch/v1/capability-model.md` — agent/capability model, reasoning dispatch, commit flow
- `docs/arch/v1/scheduler.md` — heap, inference types, routing
