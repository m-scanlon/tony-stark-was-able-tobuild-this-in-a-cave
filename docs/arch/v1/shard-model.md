# Shard Model

> **Status: Early draft. This needs significantly more design work before it can be considered canonical. The ideas here are directionally correct but many details are unresolved.**

---

## What a Shard Is

A Shard is a lightweight presence that Skyra establishes on a device. It fingerprints the device, registers its capabilities with the control plane, and listens for commands.

Not all devices are equal. Some run open operating systems where software can be freely installed. Others are locked-down appliances where nothing can be installed at all. The Shard model accounts for both.

---

## Two Modes of Deployment

### Daemon Shard

A software daemon installed directly on an open machine. This is the primary and preferred Shard type.

- Runs natively on the host OS (Mac, Linux, Windows, Android, etc.)
- Full OS access — reads hardware, software, running processes
- Fingerprints the host machine and registers its full capability profile
- Connects to the control plane over the local network
- Examples: Mac mini, dev laptop, Linux server, Raspberry Pi

### Hardware Shard (Bootstrap Tool)

A small physical compute device (HDMI stick, compute module) that you physically attach to a locked or unknown device.

- Owns its own OS — not dependent on the host
- Probes the target device via available ports (HDMI, USB) to discover what it exposes
- Attempts to install a Daemon Shard on the target device
- If installation succeeds: unplug the stick, the daemon runs independently
- If installation fails: stick stays plugged in and operates at the best available capability tier
- Not a permanent fixture — intended as a reusable bootstrap tool carried device to device

---

## Hardware Shard: Probe and Negotiate

On first boot, the Hardware Shard probes the target device across every available interface and discovers what it actually exposes. It does not assume.

### What HDMI can yield

- **EDID** — always available, no host cooperation needed. Manufacturer, model, supported resolutions, HDR capability, audio support.
- **HDMI-CEC** — device type, manufacturer, logical address, power state, supported control commands (power, volume, input switching). Many consumer devices support this.
- **ARC/eARC** — audio return channel capability.

### What USB can yield

- **Device descriptor** — vendor ID, product ID, manufacturer string, device class. Available on any USB connection.
- **Mass storage** — if exposed, browsable filesystem
- **CDC network** — if exposed, local network link, can attempt service discovery
- **CDC serial** — attempt handshake, see if anything responds
- **ADB** — if the device is Android with debugging enabled, substantial access
- **Roku developer mode** — if enabled, channel sideload is possible via USB

The probe results become the device's registered capability profile. Dynamic, not hardcoded.

---

## Capability Negotiation Chain

After probing, the Shard walks down a chain from most to least capable:

```
1. Full daemon install
   ↓ (OS is locked, cannot install)
2. Sideload app or channel (Roku dev mode, Android APK, smart TV SDK)
   ↓ (no sideload available)
3. Standard protocol control (HDMI-CEC, UPnP, ADB, device API)
   ↓ (no usable protocols)
4. EDID only — know what the device is, no control surface
   ↓ (nothing)
5. Register as unknown — device present, no control surface available
```

Whatever level it lands on becomes the registered capability profile. The control plane receives the result, not the process.

If a daemon is successfully installed, the Hardware Shard's role ends. The daemon handles its own identity and registration. The stick can be removed.

---

## Probe Layers

Probing is not limited to physical ports. A full probe runs across four layers, with each layer available depending on the Shard type and network context.

### Physical (Hardware Shard only)

| Interface | What you get |
|---|---|
| HDMI — EDID | Manufacturer, model, resolutions, HDR, audio. Always available, no cooperation needed. |
| HDMI — CEC | Device type, power state, supported control commands (power, volume, input). |
| HDMI — ARC/eARC | Audio return channel capability. |
| USB — device descriptor | Vendor ID, product ID, manufacturer string, device class. |
| USB — mass storage | Browsable filesystem if exposed. |
| USB — CDC network | Local link, enables service discovery. |
| USB — CDC serial | Attempt handshake, see if anything responds. |
| USB — ADB | Android with debugging enabled — substantial access. |
| IR | Learn remote codes via receiver. Control via blaster. Universal fallback for TVs, AV receivers, AC units. |

### Wireless local

| Interface | What you get |
|---|---|
| Bluetooth scan | BLE advertisements are passive and constant. Device name, manufacturer data, service UUIDs (reveals supported protocols), signal strength for proximity. Zero interaction required. |
| WiFi scan | Nearby SSIDs, signal strength, channel, security type. Some devices embed their type in the SSID (Sonos-xxxx, HP-Printer-xxxx). No joining required. |
| Zigbee / Z-Wave / Thread / Matter | With the right radio, direct discovery of smart home devices — lights, locks, sensors, plugs, thermostats. |

### Network (trusted networks only — see Network Trust Boundary below)

| Method | What you get |
|---|---|
| mDNS / Bonjour (passive) | Devices announce service types continuously. `_airplay._tcp`, `_googlecast._tcp`, `_hap._tcp` (HomeKit), `_raop._tcp`, `_ssh._tcp`, `_smb._tcp`. Rich, free, no scanning. |
| SSDP / UPnP (passive) | Roku, smart TVs, media servers announce via UDP multicast. Roku ECP endpoint discovered this way. |
| ARP scan | All IPs and MAC addresses on the subnet. MAC OUI prefix identifies manufacturer — cross-reference against OUI database to get probable device type before any direct communication. |
| DHCP table | If accessible via router, gives full device list with IPs, MACs, and hostnames in one query. |
| Port scan + known service probing | For non-announcing devices. Probe known ports: Roku 8060, Sonos 1400, Hue 80/443. Confirm service with handshake. |
| HTTP fingerprinting | Hit common web ports, read response headers and SSL certificate data. Many devices identify themselves in HTTP responses. |
| SSH banner | If port 22 is open, read the banner. Often contains OS and version. Non-intrusive, read-only. |
| SNMP | Routers, switches, printers, servers. Query for device info, interface stats, system description. |

### Application

Once a device is identified, probe its specific API. Roku ECP, Chromecast local API, Sonos local API, HomeKit HAP, UPnP media control. This is where actual capability registration happens.

---

## Network Trust Boundary

**A Hardware Shard on an unknown or unregistered network performs physical probes only — HDMI and USB on the target device. It does not scan the network.**

Many consumer IoT devices have no authentication on their local API — Roku ECP, Chromecast, Sonos, and others are fully controllable by any device on the same network. Running network-wide discovery on someone else's WiFi (a friend's house, a hotel, a coffee shop) would constitute unauthorized access to their devices. This is a hard boundary in the design.

Network-wide discovery runs only on networks explicitly registered as trusted by the user.

---

## Registration Wizard

The user's phone is the consent layer. No network-wide access is granted automatically — the user explicitly approves what Skyra can see and control during a guided registration flow.

Flow:

```
1. BLE provisioning — phone sends WiFi credentials to Hardware Shard
2. Shard joins network, probes target device (physical only at this point)
3. Wizard: "Found [device]. Scan for other devices on this network?"
4. User approves → network discovery runs → discovered devices presented
5. User selects what Skyra is allowed to access
6. Confirmed grants registered with control plane as part of Shard profile
7. Anything not explicitly granted = blocked by default
```

If the user is at someone else's house — no registration wizard for that network, no network discovery, physical probe of the target device only.

Access grants are stored in the Shard's capability profile and enforced by the boundary layer. The same boundary model used for agent tool access applies here — explicit allow list, everything else denied.

---

## Network Setup: WiFi Provisioning over BLE

The Hardware Shard has no WiFi credentials on first boot. It uses Bluetooth to receive them from the user's phone before connecting to the local network.

Flow:
```
1. Shard boots — no WiFi credentials
2. Advertises over BLE
3. Skyra app on phone discovers it
4. App sends WiFi credentials over BLE
5. Shard connects to local network
6. Reaches control plane
7. BLE provisioning complete
```

This is a standard IoT onboarding pattern (same approach used by Google Home, Nest, Echo). Well-supported on ESP32, which has both WiFi and Bluetooth on a single chip and ships a provisioning library for this exact flow.

---

## Full Bootstrap Flow

```
Plug in Hardware Shard
  → BLE provisioning (phone sends WiFi credentials)
  → Connects to local network
  → Reaches control plane
  → Probes target device (HDMI + USB simultaneously)
  → Walks capability negotiation chain
  → Installs daemon if possible
  → Registers capability profile with control plane
  → Unplug (if daemon installed successfully)
    OR stay plugged in (if daemon could not be installed)
```

---

## What Needs More Design

This model is directionally correct but not fully specified. Known open areas:

- **Hardware spec** — exact form factor, chip selection, port configuration. ESP32-based is the current direction but not locked.
- **Daemon distribution** — how does the Hardware Shard obtain the right daemon binary for the target architecture? Fetched from control plane? Bundled on the stick?
- **Auth and identity** — how does a newly installed daemon establish trust with the control plane? Key generation at install time is the current thinking, but the protocol is not defined.
- **Capability profile schema** — what does a registered capability profile actually look like? What fields? How does the control plane use it for routing?
- **Fallback Shard behavior** — when the stick stays plugged in permanently (locked device), how does it receive and execute commands? What is the protocol between the stick and the control plane?
- **BLE provisioning app** — Skyra mobile app is not yet designed. This is a dependency.
- **Re-provisioning** — what happens when WiFi credentials change or the device moves to a new network?
- **Shard lifecycle** — registration, deregistration, capability updates, versioning.
- **Security** — the stick has access to whatever it can probe on the target device. What are the trust boundaries? How is the stick itself authenticated before it's allowed to register a new device?
