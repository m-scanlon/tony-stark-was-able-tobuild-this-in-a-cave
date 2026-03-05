# Distributed Brain

> **Status: Basic outline only. This needs significantly more thought before any design decisions can be made. Treat this as a starting point for future exploration, not a specification.**

---

## The Core Idea

Skyra is not tied to the Mac mini. She is tied to whatever hardware best serves her current needs. When the topology of available Shards changes — a new machine comes online, a machine goes offline, resources shift — Skyra reacts and reorganizes accordingly.

The Mac mini is currently the brain because it is the most capable available node. That is not a permanent assignment. It is the current best option.

---

## Self-Organizing Infrastructure

Each Shard that comes online reports not just its static capability profile at registration, but its current resource state on an ongoing basis — CPU load, available RAM, GPU utilization, disk, network. The control plane maintains a live capability registry across all known Shards.

When a new machine joins with more RAM or processing power, Skyra detects it, evaluates it against current workload placement, and migrates services as appropriate. The infrastructure upgrades itself.

**What this requires:**

- **Continuous resource reporting** — Shards stream current state back to the control plane, not just at registration. This is a heartbeat with resource data attached.
- **Workload catalog** — every service Skyra runs is described with its requirements. Minimum RAM, GPU needed or preferred, latency constraints, colocation requirements.
- **Placement engine** — continuously matches workloads to Shards based on live resource state and workload requirements. Re-evaluates on topology changes.
- **Stability thresholds** — migration only triggers when improvement is significant and the new host has been stable long enough to trust. Prevents thrashing.
- **Migration protocol** — how services actually move. Stateless services restart trivially. Stateful services need snapshot, transfer, and replay. Databases need consistent cutover. LLM models are large file transfers. The control plane itself is the hard case.

---

## Master-Slave with Automatic Failover

The system runs a primary-replica architecture. One node is the current brain (Leader). All other capable nodes are replicas (Followers), continuously synchronized and ready to be promoted.

This maps closely to the Raft consensus algorithm — the same pattern used by etcd and Kubernetes.

### Node States

```
Leader     → current brain. Runs the control plane. Replicates state to followers.
Follower   → replica. Stays in sync. Ready to be promoted.
Candidate  → triggered when leader heartbeat is lost. Running for election.
```

### Heartbeat and Failure Detection

The Leader sends heartbeats to all Followers continuously. If a Follower does not hear a heartbeat within a timeout window, it does not wait — it becomes a Candidate and calls an election. No human intervention required.

If the Leader detects it is going down gracefully (resource exhaustion, planned shutdown), it broadcasts a distress signal to all known nodes and can recommend a successor. This is faster than waiting for timeout but the system recovers either way.

### Capability-Aware Election

Standard Raft elects the node with the most complete log. Skyra adds a capability dimension.

```
Election criteria:
  1. Log must be current (standard Raft requirement)
  2. Node must meet minimum capability threshold (RAM, CPU floor)
  3. Among eligible candidates, most capable node wins
```

A Pi cannot be elected brain even if it has the freshest log. A node below the capability floor is a Follower only.

---

## Connection to the Object Store

The existing object store is already shaped for this. Commits are immutable, append-only, and replayable. This is structurally identical to a Raft log. If commits are replicated from Leader to Followers in real time, any Follower can reconstruct full current state by replaying the log from scratch.

The foundation is already there. Replication is the missing piece.

---

## Migration Complexity by Service Type

| Service | Complexity | Notes |
|---|---|---|
| Stateless services | Low | Stop on old host, start on new host |
| Stateful services | Medium | Snapshot, transfer, replay |
| SQLite databases | Medium | Consistent snapshot, verify integrity, cutover |
| LLM models | Low-Medium | Large file transfer, then process restart |
| Vector stores | Medium | Snapshot and restore |
| Control plane itself | High | Requires leader election before cutover. Cannot turn off the brain mid-migration. |

---

## How Far This Goes

```
Level 1 — workload migration between existing registered Shards
Level 2 — model migration (run a model on a different machine when preferred host is unavailable)
Level 3 — control plane migration (brain moves to most capable available host)
Level 4 — dynamic compute acquisition (provision cloud VM when local capacity is insufficient)
```

Level 1 is the near-term target. Level 4 means Skyra provisioning her own cloud infrastructure on demand — she is no longer just managing your hardware, she is acquiring new compute as needed. That is a significantly different thing and is noted here only to understand where this leads.

---

## Build Horizon

This infrastructure layer is a future horizon — not the current build. The Mac mini is the brain, it is hardcoded, and that is correct for now.

**What "designing toward this" means in practice:**

You do not need to implement any of this to start building. You need to avoid decisions that would have to be torn out later. Specifically:

- Do not assume there is only one control plane node in the code
- Keep the object store append-only and replayable — it is already shaped correctly
- Keep capability profiles as a first-class concept from day one
- Put a thin abstraction between the control plane and "where state lives" so the backing store can be swapped

These are cheap to do now and expensive to retrofit later. None of them require implementing Raft.

---

## Infrastructure Layer: Build vs Use

The infrastructure layer — consensus, replication, health checking, cluster membership — should not be built from scratch. This is solved engineering. Two serious open source options:

### etcd

Distributed key-value store with Raft built in. Battle-tested, Go native. Used by Kubernetes as its cluster state store. Handles consensus, leader election, and watch/notify out of the box.

Run etcd as the infrastructure layer. Everything else sits on top of it. Clean separation.

### NATS with JetStream

Lighter weight than etcd, and more interesting for Skyra specifically because it solves two problems at once. NATS is a messaging system — Skyra needs a message bus for events between services anyway. JetStream adds persistence and clustering. One dependency instead of two.

**NATS is the more likely direction.** It covers the event bus that the current architecture already needs and adds clustering on top. Evaluate seriously before committing.

### The interface that matters

Whichever is chosen, the control plane and Shards should only ever talk to a thin abstraction over it. The interface is small:

```
Am I the current leader?
Store this state
Get current state
Watch for state changes
Who are the current healthy nodes and their capabilities?
```

Swap the implementation underneath without touching anything above.

---

## What Needs Much More Design

Almost everything here is directional. None of it is specified.

- **etcd vs NATS** — needs proper evaluation before any infrastructure work begins. NATS is the current preference but this is not locked.
- **Raft vs simpler alternatives** — full Raft is complex. For a small personal cluster (3-5 nodes), a simpler leader election mechanism may be sufficient. Evaluate against etcd/NATS before deciding to roll anything custom.
- **Replication protocol** — how commits stream from Leader to Followers. Frequency, consistency guarantees, conflict handling.
- **Capability floor definition** — minimum specs for a node to be eligible as brain. Will vary as Skyra's services evolve.
- **Placement engine design** — algorithm that matches workloads to nodes. How it weights RAM vs CPU vs GPU vs latency. How it handles colocation requirements.
- **Migration orchestration** — who coordinates a migration? What happens if it fails mid-transfer?
- **Split-brain handling** — if the network partitions and two nodes both think they are Leader, state can diverge. Needs a defined resolution strategy.
- **Security** — nodes accepting promotions and state from other nodes is a significant attack surface. Trust model between nodes needs definition.
- **Interaction with Shard registration** — capability registry feeds both placement and election. How are these kept consistent?
- **Level 4 scope and safety** — if Skyra can provision cloud compute, what are the guardrails? Cost limits, geographic constraints, data residency.
