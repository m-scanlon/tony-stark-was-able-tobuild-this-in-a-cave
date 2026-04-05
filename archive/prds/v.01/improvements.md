# Improvements

## Perception and History

JSON objects do not belong in perception. They are wasted space.

We need to manage perception better and actually define history, because the current code is assuming it is just past interaction, and that is not the case.

Perception should contain only the bounded working state needed for cognition, not hydrated transport payloads or command envelopes.

Current `history` in the kernel should be treated as short-term runtime context, not the real long-term history system.

The real history / retrieval / knowledge graph layer is a separate service boundary and should not be implicitly collapsed into the kernel's in-memory state.

## Suspended Chains

Suspended chains should carry a lightweight resumable contract, something like:

- prior goal
- last stable perception
- unresolved question
- suspension reason
- stale-after timestamp

## Multi-Stimulus Concurrency

We need to push multi-stimulus concurrency and see where the current design breaks.

Likely break points in the current design:

- perception is effectively a singleton, which means concurrent chains do not have clean isolation
- activeChain and activeStep are singletons, which means the runtime only has one real reasoning slot
- history is global and underspecified, so concurrent chains can read shifting context without a clear boundary
- lastUnderstanding is global, which makes cross-chain leakage likely
- the interrupt model assumes one active chain being preempted by one higher-priority stimulus
- suspended chains currently preserve too much raw chain state and not enough intentional resume state
- the thought surface assumes one active chain and one step stream
- the heap is global priority ordering, not a concurrency scheduler with fairness, quotas, or chain-local backpressure
- interaction writes go back into shared history, which can pollute unrelated concurrent reasoning paths

If we push toward multi-stimulus concurrency, we likely need:

- per-chain perception instead of one mutable runtime perception
- an actual definition of history, including what belongs in it and what does not
- a resumable contract for suspended or background chains
- a scheduler for multiple live chains instead of one active chain plus suspension
- explicit rules for shared state writes, understanding carry-forward, and chain merge or discard behavior

## Protocol Boundary

The protocol should stay at the `skyra <primitive> <args>` command boundary, and those commands should pass through the kernel.

Commands may carry hydrated JSON objects when needed, but those JSON objects should not enter perception.

Stimuli should be passed as `skyra experience <args>`.

Ingress should normalize external input into that command boundary before the kernel reasons over it.

## Model Contract

We need to lock down the contract that deals with calling the model.

Once the chain-of-thought design is in a shape we actually like, we can do the more painful exercise of defining the broader runtime contracts together.
