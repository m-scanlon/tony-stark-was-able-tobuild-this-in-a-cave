# Adapter Writer

The being that closes the self-extension loop. When Skyra needs a new surface she cannot reach, she fires at the adapter-writer. It generates the adapter program, writes it to disk, and emits a grow directive. The router picks it up. The new being is live. No human touched the code.

## What It Is

A cognitive being. It reasons through inference. Its sole purpose is to generate adapter programs that implement the wire format.

It has a process surface like any other cognitive being — backed by an inference adapter pointed at a model. What makes it different is not its execution surface but its purpose: it produces programs, not protocol strings alone.

## What Skyra Says To It

```
skyra adapter-writer ~name <being-name> ~connects-to <description of the surface> ~purpose <what the being will do> | I need to reach X
```

The description is freeform. Skyra describes what she needs to connect to. The adapter-writer reasons about what the adapter program needs to do and generates it.

## What It Produces

Two things come back as protocol strings:

**1. The adapter program on disk.**

The adapter-writer writes a program to a known adapters directory. The program implements the wire format — reads present from stdin until `---`, does its work, writes protocol strings to stdout, writes `---`. The language does not matter. The wire format is the only contract.

**2. A grow directive.**

```
skyra grow ~name <being-name> ~surface process ~command <path-to-adapter> ~identity <identity> ~purpose <purpose> | adapter-writer: new surface registered
```

This is a normal protocol string emitted back into the runtime. grow processes it. The router reads the surface declaration, spawns the adapter, holds the handles. The being is reachable from that moment.

## Why This Works

The wire format is stable. It does not change. Every adapter — human-written or generated — speaks the same contract: present in, protocol strings out, `---` as terminator. The adapter-writer knows this contract and generates programs that implement it. As long as the wire format holds, every generated adapter works without modification to the runtime.

grow is already the control surface for beings. The adapter-writer does not need a special registration mechanism — it emits a grow directive the same way any being would. The runtime does not know the being was generated. It just sees a grow directive and does what it always does.

## The Loop Closing

```
Skyra decides she needs a new surface
  → fires at adapter-writer
  → adapter-writer generates program, writes to disk
  → adapter-writer emits grow directive
  → grow registers the being
  → router spawns the adapter
  → new being is reachable
  → Skyra fires at it
```

Every step uses existing machinery. No new mechanism. The adapter-writer is just a being that happens to produce other beings as output.

## What The Adapter-Writer Needs To Know

It carries in its identity and purpose the wire format specification and the adapters directory path. These are seeded in the genome. When Skyra opens an exchange with the adapter-writer, its present includes this grounding — it knows the contract it is generating against.

The callable language on the adapter-writer's channel makes the required fields explicit:

```
~name <being-name> ~connects-to <surface description> ~purpose <what the being will do>
```

Inference cannot omit these. If they are missing the runtime rejects the expression before it reaches the adapter-writer.

## Dependency On Wire Format Stability

The adapter-writer generates programs against the wire format. If the wire format changes, every generated adapter breaks. The wire format is therefore a hard invariant of the runtime — it does not change without regenerating all adapters.

This is the one constraint that must hold for self-extension to work.

## Open Questions

- What language does the adapter-writer generate programs in? Go requires compilation. Python or shell script runs immediately. The simplest path to Skyra running a generated adapter fast is a scripting language.
- Where is the adapters directory? Declared in the genome, or fixed by the runtime?
- Does the adapter-writer test the generated adapter before emitting the grow directive — run it, send a test present, verify it responds correctly?
- What happens if the generated adapter is wrong — it writes invalid protocol strings or crashes immediately? Does the adapter-writer get the error back and retry?
- Can the adapter-writer write adapters that are themselves cognitive — backed by inference — or only non-cognitive transducers?
