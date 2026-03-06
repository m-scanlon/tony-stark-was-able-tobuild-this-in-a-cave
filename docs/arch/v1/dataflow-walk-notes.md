# Data Model Walk — Misc Notes

Decisions and observations captured while walking the pipeline service by service. Not formal design docs — just things worth remembering.

---

## Event Register (Voice Shard only, for now)

The Voice Shard maintains an event register — a hash table keyed by `turn_id` — to track events that have been sent but not yet ACKed.

- Send event → insert entry
- ACK received → lookup by `turn_id`, pop
- Retry loop → scan for entries past `next_attempt_at`
- Durable (SQLite) so Pi reboots don't lose in-flight events

**The Brain Shard does not need an equivalent structure for ACKing.** The ACK leg is synchronous — Brain Shard receives, writes to SQLite inbox, sends ACK back on the same connection. No tracking structure needed to facilitate that. The Brain Shard's SQLite inbox is for downstream pipeline consumption, not ACK mechanics.

**GPU Shard probably doesn't need an event register.** In v1 it's a dumb inference endpoint — request/response, no async event tracking. Revisit when we walk the GPU Shard's role. If it ever handles async or multi-step work independently, it'd need one.

**Future: Shard bootstrap package.** Brain Shard knows each Shard's capability profile from registration. It could provision each Shard with a tailored bootstrap package — Voice Shard gets event register + retry loop + voice config, GPU Shard gets inference config only. Package contents derived from registered capabilities. Capture properly when we get to Shard provisioning.

---

## ACK uses turn_id, not event_id

`event_id` is internal to the Brain Shard — generated on ingress, never crosses the wire to Voice Shard. ACK references `turn_id` only. Voice Shard clears its event register by `turn_id`.

---
