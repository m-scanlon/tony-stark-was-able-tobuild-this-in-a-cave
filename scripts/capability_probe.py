#!/usr/bin/env python3
"""
Capability probe v0.

Purpose:
- Reveal local machine capabilities without needing a daemon protocol.
- Optionally probe an explicitly named target host for remotely exposed
  capabilities without doing a subnet-wide scan.

This script is intentionally conservative:
- it never scans an entire subnet
- it only probes the local host plus targets you explicitly name
- it uses standard-library networking only

Examples:
  python3 scripts/capability_probe.py
  python3 scripts/capability_probe.py --target 192.168.1.50
  python3 scripts/capability_probe.py --target roku.local --ports 80,443,8060
  python3 scripts/capability_probe.py --target 192.168.1.50 --output probe.json
"""

from __future__ import annotations

import argparse
import datetime as dt
import http.client
import json
import os
import platform
import shutil
import socket
import ssl
import subprocess
import sys
from typing import Any


DEFAULT_PORTS = [
    22,    # ssh
    53,    # dns
    80,    # http
    443,   # https
    5555,  # adb over tcp
    7000,  # airplay-ish
    7001,  # airplay tls-ish
    8008,  # chromecast-ish
    8009,  # chromecast-ish
    8060,  # roku ecp
    1400,  # sonos
]


def now_iso() -> str:
    return dt.datetime.now(dt.timezone.utc).isoformat().replace("+00:00", "Z")


def run_command(cmd: list[str]) -> str | None:
    try:
        result = subprocess.run(
            cmd,
            check=False,
            capture_output=True,
            text=True,
        )
    except FileNotFoundError:
        return None

    if result.returncode != 0:
        return None

    return result.stdout


def parse_json_command(cmd: list[str]) -> dict[str, Any] | None:
    stdout = run_command(cmd)
    if not stdout:
        return None

    try:
        return json.loads(stdout)
    except json.JSONDecodeError:
        return None


def parse_memory_bytes(raw: str | None) -> int | None:
    if not raw:
        return None

    parts = raw.split()
    if len(parts) != 2:
        return None

    try:
        value = float(parts[0])
    except ValueError:
        return None

    unit = parts[1].upper()
    scale = {
        "KB": 1024,
        "MB": 1024**2,
        "GB": 1024**3,
        "TB": 1024**4,
    }.get(unit)
    if scale is None:
        return None

    return int(value * scale)


def system_profiler_payload(types: list[str]) -> dict[str, Any]:
    if platform.system() != "Darwin":
        return {}

    payload = parse_json_command(["system_profiler", "-json", "-detailLevel", "mini", *types])
    return payload or {}


def add_capability(
    capabilities: list[dict[str, Any]],
    *,
    name: str,
    status: str,
    verification: str,
    commands: list[str] | None = None,
    constraints: list[str] | None = None,
    evidence: list[dict[str, Any]] | None = None,
) -> None:
    capabilities.append(
        {
            "name": name,
            "status": status,
            "verification": verification,
            "commands": commands or [],
            "constraints": constraints or [],
            "evidence": evidence or [],
        }
    )


def local_fingerprint() -> dict[str, Any]:
    disk_total, disk_used, disk_free = shutil.disk_usage("/")
    host = {
        "hostname": socket.gethostname(),
        "platform": platform.system(),
        "platform_release": platform.release(),
        "platform_version": platform.version(),
        "architecture": platform.machine(),
        "python_version": platform.python_version(),
        "logical_cpu_count": os.cpu_count(),
        "disk_total_bytes": disk_total,
        "disk_used_bytes": disk_used,
        "disk_free_bytes": disk_free,
    }

    if platform.system() == "Darwin":
        hw = system_profiler_payload(["SPHardwareDataType"]).get("SPHardwareDataType", [])
        if hw:
            item = hw[0]
            host["machine_name"] = item.get("machine_name")
            host["machine_model"] = item.get("machine_model")
            host["chip_type"] = item.get("chip_type")
            host["physical_memory_bytes"] = parse_memory_bytes(item.get("physical_memory"))

    return host


def local_capabilities() -> list[dict[str, Any]]:
    capabilities: list[dict[str, Any]] = []
    fingerprint = local_fingerprint()

    add_capability(
        capabilities,
        name="local_compute",
        status="verified",
        verification="host_os_introspection",
        commands=["compute.execute"],
        evidence=[
            {"kind": "platform", "value": fingerprint.get("platform")},
            {"kind": "chip_type", "value": fingerprint.get("chip_type")},
            {"kind": "logical_cpu_count", "value": fingerprint.get("logical_cpu_count")},
            {"kind": "physical_memory_bytes", "value": fingerprint.get("physical_memory_bytes")},
        ],
    )

    if fingerprint.get("disk_total_bytes"):
        add_capability(
            capabilities,
            name="local_storage",
            status="verified",
            verification="disk_usage_root_fs",
            commands=["storage.read", "storage.write"],
            evidence=[
                {"kind": "disk_total_bytes", "value": fingerprint.get("disk_total_bytes")},
                {"kind": "disk_free_bytes", "value": fingerprint.get("disk_free_bytes")},
            ],
        )

    if platform.system() == "Darwin":
        payload = system_profiler_payload(
            [
                "SPDisplaysDataType",
                "SPCameraDataType",
                "SPBluetoothDataType",
                "SPUSBDataType",
            ]
        )

        displays = payload.get("SPDisplaysDataType", [])
        if displays:
            evidence = []
            for display in displays:
                evidence.append(
                    {
                        "kind": "gpu",
                        "model": display.get("sppci_model") or display.get("_name"),
                        "metal": display.get("spdisplays_metal"),
                        "cores": display.get("sppci_cores"),
                    }
                )
            add_capability(
                capabilities,
                name="display_output",
                status="verified",
                verification="system_profiler SPDisplaysDataType",
                commands=["display.render"],
                evidence=evidence,
            )

            metal_supported = any(d.get("spdisplays_metal") == "spdisplays_supported" for d in displays)
            if metal_supported:
                add_capability(
                    capabilities,
                    name="gpu_acceleration",
                    status="verified",
                    verification="system_profiler SPDisplaysDataType",
                    commands=["compute.accelerated"],
                    evidence=evidence,
                )

        cameras = payload.get("SPCameraDataType", [])
        if cameras:
            add_capability(
                capabilities,
                name="camera_input",
                status="verified",
                verification="system_profiler SPCameraDataType",
                commands=["camera.capture"],
                evidence=[{"kind": "camera", "value": item.get("_name", "camera")} for item in cameras],
            )

        bluetooth = payload.get("SPBluetoothDataType", [])
        if bluetooth and any(item for item in bluetooth if item):
            add_capability(
                capabilities,
                name="bluetooth",
                status="verified",
                verification="system_profiler SPBluetoothDataType",
                commands=["bluetooth.scan", "bluetooth.connect"],
                evidence=bluetooth,
            )

        usb = payload.get("SPUSBDataType", [])
        if usb:
            add_capability(
                capabilities,
                name="usb_host",
                status="verified",
                verification="system_profiler SPUSBDataType",
                commands=["usb.enumerate"],
                evidence=[
                    {
                        "kind": "usb_bus",
                        "name": item.get("_name"),
                        "host_controller": item.get("host_controller"),
                    }
                    for item in usb
                ],
            )

    return capabilities


def tcp_connect(host: str, port: int, timeout: float) -> tuple[bool, str | None]:
    try:
        with socket.create_connection((host, port), timeout=timeout):
            return True, None
    except OSError as exc:
        return False, str(exc)


def http_probe(host: str, port: int, use_tls: bool, timeout: float) -> dict[str, Any] | None:
    try:
        if use_tls:
            conn = http.client.HTTPSConnection(
                host,
                port=port,
                timeout=timeout,
                context=ssl._create_unverified_context(),
            )
        else:
            conn = http.client.HTTPConnection(host, port=port, timeout=timeout)

        conn.request("GET", "/", headers={"User-Agent": "skyra-capability-probe/0.1"})
        response = conn.getresponse()
        headers = {k.lower(): v for k, v in response.getheaders()}
        body = response.read(256)
        conn.close()
        return {
            "status": response.status,
            "reason": response.reason,
            "server": headers.get("server"),
            "content_type": headers.get("content-type"),
            "body_prefix": body.decode("utf-8", errors="replace"),
        }
    except Exception as exc:  # noqa: BLE001 - probe should degrade softly
        return {"error": str(exc)}


def classify_open_port(port: int, http_info: dict[str, Any] | None) -> tuple[str, list[str]]:
    if port == 22:
        return "ssh_endpoint", ["remote.shell"]
    if port in (80, 443):
        return "http_api", ["http.request"]
    if port == 5555:
        return "adb_tcp", ["android.debug"]
    if port in (7000, 7001):
        return "airplay_like_endpoint", ["media.cast"]
    if port in (8008, 8009):
        return "chromecast_like_endpoint", ["media.cast"]
    if port == 8060:
        return "roku_ecp_endpoint", ["device.control", "media.launch"]
    if port == 1400:
        return "sonos_http_endpoint", ["audio.control"]

    if http_info and "status" in http_info:
        return "http_api", ["http.request"]

    return f"tcp_port_{port}", []


def probe_target(target: str, ports: list[int], timeout: float) -> dict[str, Any]:
    result: dict[str, Any] = {
        "target": target,
        "resolved_addresses": [],
        "capabilities": [],
        "open_ports": [],
        "probe_errors": [],
    }

    try:
        addrinfo = socket.getaddrinfo(target, None, type=socket.SOCK_STREAM)
        addresses = sorted({item[4][0] for item in addrinfo})
        result["resolved_addresses"] = addresses
    except socket.gaierror as exc:
        result["probe_errors"].append({"kind": "dns", "error": str(exc)})
        return result

    for port in ports:
        is_open, error = tcp_connect(target, port, timeout)
        if not is_open:
            continue

        entry: dict[str, Any] = {"port": port}
        http_info = None
        if port in (80, 443, 1400, 8060, 8008, 8009, 7000, 7001):
            http_info = http_probe(target, port, use_tls=(port == 443), timeout=timeout)
            if http_info:
                entry["http"] = http_info

        name, commands = classify_open_port(port, http_info)
        result["open_ports"].append(entry)
        add_capability(
            result["capabilities"],
            name=name,
            status="verified",
            verification=f"tcp_connect:{port}",
            commands=commands,
            evidence=[entry],
            constraints=[
                "remote capability only proves exposed surface",
                "no daemon-level local execution assumed",
            ],
        )

    if not result["capabilities"]:
        add_capability(
            result["capabilities"],
            name="network_presence_only",
            status="partial",
            verification="dns_resolution_only",
            constraints=["target resolved but no tested service ports answered"],
            evidence=[{"kind": "resolved_addresses", "value": result["resolved_addresses"]}],
        )

    return result


def parse_ports(raw: str | None) -> list[int]:
    if not raw:
        return list(DEFAULT_PORTS)

    ports: list[int] = []
    for part in raw.split(","):
        part = part.strip()
        if not part:
            continue
        port = int(part)
        if port < 1 or port > 65535:
            raise ValueError(f"invalid port: {port}")
        ports.append(port)
    return ports


def build_report(targets: list[str], ports: list[int], timeout: float) -> dict[str, Any]:
    report: dict[str, Any] = {
        "probe_kind": "capability_probe_v0",
        "generated_at": now_iso(),
        "local_subject": {
            "kind": "local_host",
            "fingerprint": local_fingerprint(),
            "capabilities": local_capabilities(),
        },
        "remote_subjects": [],
    }

    for target in targets:
        report["remote_subjects"].append(probe_target(target, ports, timeout))

    return report


def main() -> int:
    parser = argparse.ArgumentParser(description="Reveal capability surfaces on the local host and explicit remote targets.")
    parser.add_argument(
        "--target",
        action="append",
        default=[],
        help="Explicit hostname or IP to probe for remotely exposed capabilities. May be used multiple times.",
    )
    parser.add_argument(
        "--ports",
        help="Comma-separated port list for remote target probing. Default is a small capability-oriented set.",
    )
    parser.add_argument(
        "--timeout",
        type=float,
        default=0.75,
        help="Socket timeout in seconds for remote probes. Default: 0.75",
    )
    parser.add_argument(
        "--output",
        help="Optional path to write the JSON report.",
    )
    args = parser.parse_args()

    try:
        ports = parse_ports(args.ports)
    except ValueError as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 2

    report = build_report(args.target, ports, args.timeout)
    payload = json.dumps(report, indent=2, sort_keys=False)

    if args.output:
        with open(args.output, "w", encoding="utf-8") as f:
            f.write(payload)
            f.write("\n")

    print(payload)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
