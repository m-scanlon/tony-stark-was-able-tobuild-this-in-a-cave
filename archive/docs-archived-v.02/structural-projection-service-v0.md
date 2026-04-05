# Structural Projection Service v0 (Superseded)

This document is historical and is not canonical for this version.

It described an older direction where interaction, recall, and runtime artifacts were projected into a scored episode-side field before recall.

The current direction is:

- recall is driven by heavy inference over current episode context
- heavy inference emits bounded recall request stimulus when recall is needed
- recall results are written into `episode.recall`
- no separate structural projection service is required by the active recall contract

Historical value that still survives here:

- interaction, prior recall, and runtime artifacts may still matter as inference inputs
- canonical structure and `anchor_set` overlap still matter
- bounded preprocessing may still exist later, but it is not the current recall contract

This file should be treated as superseded context rather than active canon.
