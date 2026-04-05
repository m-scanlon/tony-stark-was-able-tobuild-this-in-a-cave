# Episode Field Scoring (Superseded)

This document is historical and is not canonical for this version.

It described an older direction where recall depended on a scored episode-side field inside the episode.

The current direction is:

- recall is a contract
- recall is driven by heavy inference
- heavy inference emits bounded recall request stimulus when recall is needed
- retained artifacts are retrieved through `anchor_set` overlap
- admitted results are written into `episode.recall`
- no separate scored field object is required in the active contract set

The historical value that still survives here is narrower:

- bounded ranking still matters
- structural overlap still matters
- relationship-aware retrieval still matters
- multi-hop recall should remain deferred unless explicitly bounded

This file should be read only as superseded context for how the older scoring-first model evolved.
