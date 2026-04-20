# Exchange Hashmap — Thread As Key

## The Change

The relationship substrate changes from a stack to a hashmap keyed by thread.

**Before:**
```
HashMap<being_id, HashMap<peer_being_id, Stack<Exchange>>>
```

**After:**
```
HashMap<being_id, HashMap<peer_being_id, HashMap<thread_id, Exchange>>>
```

## Why

A stack is ordered by time. Exchanges are not a queue — they are distinguished by what they are for. Thread is the natural key. Multiple exchanges with the same peer are now addressable by thread. No stack peeking to find the open one. No ambiguity about which exchange is active. You look up by thread.

## What thread_id Is

The thread_id is the `~about` value from `start-exchange` — what the opener is trying to resolve. It is set by the opener and travels with the exchange as its identity. Private to the opener. It keys the exchange on the opener's side.

## Parallel Exchanges

Multiple exchanges with the same peer are now structurally distinguishable. Each has its own thread. Each is addressable independently. This is how parallel threads per peer become possible without special threading machinery — the thread_id is the thread.
