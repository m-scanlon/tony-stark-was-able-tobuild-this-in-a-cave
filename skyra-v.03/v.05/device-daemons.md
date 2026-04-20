# Device Daemons

Skyra runs on one device. Everything else is a surface.

## The Brain Device

The brain device is where the runtime lives. Beings, exchanges, threads, the genome, all world state — authoritative, singular, on one machine. This is not a distributed system. The runtime does not replicate. It does not shard. There is one world.

## Daemons

Every other device — phone, laptop, watch, tablet — runs a daemon. The daemon is not a runtime. It has no beings. It has no exchange state. It is two things only:

- **Surface registry** — what execution surfaces this device offers. Screen, microphone, camera, speaker, GPS, accelerometer. Whatever the device has.
- **Router** — receives signals from the brain device, routes to the right surface, sends responses back.

The daemon is lightweight by design. It does not think. It does not hold state. It is a relay between the brain device and the physical world.

## The Inter-Device Surface

On the brain device, each daemon appears as an adapter. Same wire format — present in, protocol strings out, `---` as terminator. The adapter communicates with the daemon over the network instead of stdin/stdout. The router does not know the difference.

```
brain device runtime
  → inter-device adapter (stdin/stdout to runtime)
    → network connection
      → daemon on remote device
        → device surfaces
```

From the runtime's perspective the phone is just a being with a process surface. The inter-device adapter handles the network. The daemon handles the device.

## Surface Registration

When a daemon comes online it registers its available surfaces with the brain device. This flows through the inter-device adapter as grow directives — new beings seeded into the world representing that device's surfaces. When a daemon goes offline those beings become unreachable. The router marks their adapters stopped.

New device comes online — new surfaces appear. Device goes offline — surfaces disappear. The world reflects the current topology automatically.

## Self-Extension Across Devices

When Skyra writes a new adapter targeting a device capability, the adapter-writer generates the program, the daemon on the target device registers it, and the brain device can route to it. Self-extension is not limited to the brain device. It reaches every device in the network.

## What The Daemon Does Not Do

The daemon does not reason. It does not hold exchange state. It does not make decisions. It is a surface router only. All cognition stays on the brain device. The daemon is as dumb as possible by design — the less it does, the less that can go wrong at the edge.

## Open Questions

- How does the brain device discover daemons on the network — broadcast, known addresses, a registry service?
- What is the transport between the inter-device adapter and the daemon — TCP, WebSocket, something else?
- How does the inter-device adapter handle a daemon going offline mid-exchange — buffer, drop, or surface the disconnection to the being?
- What trust model governs which daemons the brain device accepts — any device on the local network, or explicitly registered devices only?
- Does the daemon authenticate to the brain device, or does the brain device authenticate to the daemon, or both?
- When multiple daemons are online simultaneously, does Skyra route to the most recently active device or does the being decide?
