# Data Resilience

> **Not MVP.** These skills are post-MVP. The architecture supports them — the shard model is the distribution layer — but they are not required for initial boot.

---

## The Problem

The system is personal. It runs on consumer hardware. Disks fill. RAM runs out. Machines die. The system needs a defined response to each failure mode — not a crash, not data loss, not silence.

Data integrity is the first priority. The system must never lose committed data. Observational data is expendable. The committed layer is never evicted, never dropped, never lost.

---

## Degradation Ladder

When the system is under pressure, it descends this ladder in order. Each level is a skill. Each transition is a kernel event.

```
1. RAM pressure
   → drop to smaller model
   → evict observational nodes (never committed)

2. Disk pressure, shard available
   → migrate_data → partition moves to target shard
   → shard_rebalance → Redis updated → kernel routes transparently

3. Disk pressure, no shard available
   → emergency_offload → committed partition encrypted → pushed to cloud
   → system continues lean

4. No cloud configured
   → graceful drain → complete jobs in flight → stop accepting new ones
   → notify user
```

---

## System Health Monitoring

**`system_health`** — primitive skill. Fired by cron on a schedule. Reads RAM, disk, graph size, Redis memory usage. Emits events into the kernel when thresholds are crossed. The user only sees it when action is required.

```
system_health fires (cron)
  → reads: RAM available, disk available, graph size, Redis memory
  → threshold crossed → emit event → kernel → surface to user if action needed
  → no threshold → silent
```

Thresholds are relative, not absolute. 80% disk on a 256GB drive is different from 80% on a 512GB drive. The skill reasons about headroom, not raw numbers.

---

## Primitive Skills

### `system_health`

Monitors system resource state. Emits events when thresholds are crossed. Fired by cron. Silent when healthy.

### `migrate_data`

Moves a graph partition from the current shard to a target shard.

```
trigger: system_health emits disk pressure event + shard available
  → surface to user: "Graph is 80% of disk. [Shard] has 2TB available. Migrate?"
  → user approves
  → copy partition to target shard
  → verify integrity (node count, edge count, checksum)
  → only release source after verification passes
  → shard_rebalance runs
```

Data integrity constraint: migration does not release the source partition until verification passes. A failed migration leaves the source intact.

### `shard_rebalance`

Updates Redis after a migration. The kernel routes to the new shard owner transparently. No job interruption.

```
trigger: migrate_data completes successfully
  → update Redis: partition ownership → target shard
  → kernel picks up new routing on next dispatch
  → old shard releases partition
```

### `emergency_offload`

Last resort when disk pressure is critical and no shard is available.

**Constraints:**
- Committed layer only — observational data is not offloaded
- Encrypted with the user's keypair before leaving the device
- Cloud provider sees blobs, nothing else
- User configures the endpoint — S3, Backblaze, or equivalent. Skyra does not pick.
- Cloud is temporary. Data belongs on-device.

```
trigger: disk pressure + no shard available
  → surface to user: "No shard available. Offload to [endpoint]?"
  → user approves
  → encrypt committed partition (user keypair)
  → push to user-configured cloud endpoint
  → record offload manifest locally (what was offloaded, where, when)
  → system continues with remaining local data
```

### `restore_from_offload`

Pulls offloaded data back when capacity is restored.

```
trigger: new shard comes online with sufficient capacity
       OR local disk pressure resolved
  → surface to user: "Capacity available. Restore offloaded data?"
  → user approves
  → pull from cloud endpoint
  → decrypt locally (user keypair)
  → reintegrate into graph
  → verify integrity
  → clear offload manifest
```

---

## Invariants

- **Committed layer is never evicted.** Only observational nodes are candidates for dropping under RAM pressure.
- **Migration verifies before releasing.** Source data is not dropped until destination integrity is confirmed.
- **Cloud is a pressure valve, not a home.** Offloaded data is pulled back at the first opportunity.
- **User approves every transition.** No data moves without user confirmation.
- **User owns the keys.** Encryption happens on-device before any data leaves.

---

## Related

- `docs/arch/v1/memory-structure.md` — committed vs observational layers
- `docs/arch/v1/kernel.md` — primitive skills, cron service
- `docs/arch/v1/shard/shard-registration.md` — shard capability registration
