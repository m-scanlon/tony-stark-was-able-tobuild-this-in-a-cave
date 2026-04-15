# Impulse Routing Flag

Today all impulses write to the peer exchange. There's no way for a being to route an impulse to its own self-exchange instead.

A flag on the impulse — something like `~self` — tells the kernel to write to the self-exchange rather than the peer exchange.

## Why this matters

A being finishes work that satisfies a relational debt. Before paying it outward, it may want to hold the result in its own workspace — do more work on it, sit with it, decide if it's actually ready. The `~self` flag makes that possible without the being having to fire a separate self-directed signal through a roundabout path.

The debt to the peer stays open until the being explicitly fires back to them. Routing to self is not repayment. It's staging.

## Relationship to relational debt

The self-exchange is the workspace. The peer exchange is where obligations live. A being can move a finished artifact into its own workspace before deciding to pay the debt. The two are distinct destinations and the flag makes that distinction explicit at the kernel level.

## Open questions

- Exact flag syntax — `~self`, `~hold`, something else?
- Can a being fire to both destinations simultaneously — peer exchange and self exchange in one impulse — or does that require two separate impulses?
- Does the kernel enforce that a `~self` impulse is still a valid protocol string with a target name, or does the target become implicit (always self)?
- How does a staged artifact in the self-exchange surface in the being's present when it later fires outward?
