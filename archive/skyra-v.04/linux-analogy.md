# Skyra, through the Linux lens

Skyra's runtime has the same kinematics as a small Unix. Different substrate, same shape. Worth naming the correspondences because they're load-bearing — the design that works in one world works for similar reasons in the other.

## The map

| Skyra | Linux |
|---|---|
| Thread | Process |
| Being | Executable |
| Operator | Pipe |
| Medium | Syscall |
| Relation | Bytes in the pipe |
| Protocol (`skyra <operator> <args> \| reason`) | Wire format the bytes follow |
| EntityMap | Process table / FD table |
| `grow` | `exec` |
| `start-thread` | `fork` |
| `end-thread` | `exit` returning to parent |
| Impress | Parse incoming from stdin/network |
| DerivePresent | Format outgoing to stdout |

## The primitives, paired

**Thread ≈ Process.** A thread is an active execution context with accumulated state (relationships, exchanges, metadata). It has an ID, a lifetime, it can be closed. Multiple threads can exist concurrently. Each holds a group of beings in active relation with each other — like a Unix process holds a memory space with running code.

**Being ≈ Executable.** A being is the thing that gets run — an identity, a purpose, a medium, a set of operators it can invoke. A being isn't exclusively bound to any one thread; the same being can participate in many threads at once, just as `python3` the executable can be loaded into many concurrent processes.

**Operator ≈ Pipe.** Operators carry relations between beings the way pipes carry bytes between processes. Continue-thread wires one being's output to another being's input. Start-thread creates a new conduit. End-thread closes the conduit and returns control. Pipes compose in shell (`ls | grep | wc`); operators compose in Skyra via recursion across multiple continue-thread calls.

**Medium ≈ Syscall.** A medium is how a being reaches outside the runtime. Inference mediates through an HTTP call to the model API. CLI mediates through stdin/stdout syscalls. Shell mediates through fork+exec. Every medium is, at its root, a bet that some external capability is addressable and returns a string. Syscalls are the seam between a Linux process and the kernel; mediums are the seam between a being and the world.

**Relation ≈ IPC packet.** The data flowing through the pipe. In Unix, bytes. In Skyra, a structured `{ID, Origin, ThreadID, Impulse}` encoded as a human-readable protocol line. Both are the unit of transit — what gets carried, not what carries it.

**Protocol ≈ wire format.** Unix pipes are untyped bytes; the meaning of those bytes is agreed out-of-band (text, JSON, binary). Skyra's wire format is fixed — `skyra <operator> <args> | reason` — which lets any language compile a being that speaks it.

**EntityMap ≈ process table + FD table.** The kernel tracks every running process by PID and every open handle by FD. The EntityMap tracks every reachable entity by name. Routing is a map lookup in both.

## Why the shape rhymes

Unix got where it is because a few small ideas compose well:

1. **One uniform interface.** `read(fd, buf, n)` works on files, sockets, pipes, devices. One abstraction, many substrates.
2. **Composition via conduit.** Pipes turn independent programs into chains. Each program is small; the system is big.
3. **Everything addressable.** FDs, PIDs, paths — every primitive has a name the kernel can resolve.
4. **Kernel as mediator.** Processes don't touch each other; they go through the kernel.

Skyra mirrors each:

1. **`Entity` interface.** Everything in the world implements `Relate(r) Entity`. Beings, operators, invariants, threads — all addressed the same way.
2. **Continue-thread as conduit.** Beings are small — just identity + medium. The system gets large by chaining them through continue-thread recursion.
3. **Everything in the EntityMap.** Addressable by name. String lookup at dispatch time, just like path lookup at open time.
4. **Operator as mediator.** Beings don't touch each other. They go through continue-thread, which handles routing, state, recording.

## Where it diverges

Not every primitive has a clean match.

**No shared memory.** Linux processes can `mmap` shared regions and mutate together. Skyra beings have no shared state; everything happens via relations. Closer to Erlang than Unix here.

**Beings are persistent.** A process dies when its program exits; its memory goes. A being persists as long as it's in the EntityMap. Closer to actor systems or object systems than processes.

**The wire format is structured.** Pipes are byte streams with no built-in meaning; Skyra's protocol has required shape. More like a typed RPC than a raw pipe.

**Mediums are typed surfaces, not numbered calls.** Linux has ~300 syscalls, numerical. Skyra has N mediums, each a named function. Closer to a capability-based system than Unix.

**Threads ≠ Linux threads.** Confusingly, a "thread" in Skyra corresponds to a Linux *process*, not a Linux *thread* (which is a unit of parallel execution inside a process). The naming is inherited from conversational "threads," not OS threads.

## Why it matters

The Linux model is proven — multi-decade, huge scale, minimal concepts composing without limit. If your AI-runtime kinematics rhyme with Linux's, you inherit some of the same guarantees:

- **Composition is cheap.** Adding a new being is like adding a new program — no orchestration layer needed.
- **The protocol is the contract.** Any tool that produces a protocol line works, in any language. Like any Unix program that reads/writes bytes is a first-class citizen.
- **The OS doesn't know about the app.** Linux doesn't care if you're running a web server or a text editor. Skyra's runtime doesn't care if a being is an LLM, a human, or a compiled fetcher. The substrate is generic.

This isn't a metaphor — it's a structural claim. If the shape works, the runtime has the same expansion properties Unix has: a universe you can keep building in without redesigning the core.
