# UserSpace and SystemSpace

## The Idea

The runtime may eventually have two distinct spaces — UserSpace and SystemSpace. Not fully thought through yet. Worth holding.

## The Signal

HTTP adapters receive transport-level responses (200 OK) alongside content responses. The 200 OK is not signal — it is plumbing. Something has to filter it before it reaches the runtime.

That distinction — plumbing vs signal — might be the first concrete expression of the UserSpace/SystemSpace boundary. SystemSpace handles transport-level concerns. UserSpace carries actual signal.

## Pipes Can Diverge

Unix pipes are a first-class primitive. One process stdout can feed into another process stdin. Chains are possible:

```
adapter | system-filter | runtime
```

The system filter strips transport responses. The runtime only sees signal. Everything upstream is invisible to it.

This means the UserSpace/SystemSpace split could live in the pipe chain — no runtime changes required. The runtime just reads the end of the chain.

## Status

Not designed. Not in scope for v.04. Worth returning to when the adapter layer is built and the boundary concerns become concrete.
