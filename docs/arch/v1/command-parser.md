# Command Parser

## Syntax

```
skyra <tool> [args]
```

One prefix. One tool. Args. That's the entire protocol. Every shard, every direction.

```
skyra reply "You hit 4 workouts this week"
skyra fan_out -gym -home "cancel gym and turn off lights"
skyra report "gym session cancelled"
skyra check_nginx
skyra log_workout --type=run --duration=30
skyra stream --token="nginx" --valence=-0.5 --arousal=0.8
skyra ack --turn=turn_8f4c --status=stored
```

## Resolution

The API Gateway receives the command and resolves the tool against Redis:

```
1. command arrives at Ingress
2. Redis check: does this skill exist? is this shard authorized?
3. No  → rejected
4. Yes → Redis returns full skill definition
5. command args + full skill → heap as a job
6. kernel router reads skill contract → routes to capable shard
```

The shard reasons about which skill to invoke. The API Gateway validates. The kernel executes.

## Routing

All routing is driven by the skill's contract — compute requirements declared in the skill definition. The kernel router reads this and dispatches to the right shard. No hardcoded routing logic.

```
skill contract: { compute: "deep_reasoning" }  → routes to GPU shard
skill contract: { compute: "voice" }            → routes to Voice Shard
skill contract: { compute: "control_plane" }    → routes to Brain Shard
```

## Related

- `docs/arch/v1/kernel.md` — kernel router, execution model
- `docs/arch/v1/api-gateway.md` — Ingress validation, job envelope assembly
- `docs/arch/v1/shard-communication.md` — unified protocol, shard primitives
