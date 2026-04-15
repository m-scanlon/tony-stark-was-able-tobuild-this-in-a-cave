# Parallel Exchanges Per Peer

Deep question. Not sure if this is the right direction yet.

If multiple external conversations are happening simultaneously, the internal network might need to run multiple parallel threads. Prefrontal and strategy could have three open exchanges at once — one per external conversation. The current stack structure assumes one thread per peer at a time. That breaks under this model.

The exchange ID might become the key rather than just the peer name. Each internal exchange traces back to the external exchange that spawned it, keeping threads distinct all the way through.

But this is a significant structural change and it's not clear the complexity is worth it yet. Leave open.
