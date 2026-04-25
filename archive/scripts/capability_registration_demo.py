#!/usr/bin/env python3
"""
Capability registration demo using Gemini.

This is intentionally small:
- load a local probe report or generate one on the fly
- send the report to Gemini
- receive either:
  - a structured capability registration proposal
  - or raw protocol command lines for a smoke test

The API key is read from .env or the current shell environment.
Expected env vars:
- GEMINI_API_KEY
- optional GEMINI_MODEL
"""

from __future__ import annotations

import argparse
import json
import os
import pathlib
import re
import ssl
import sys
import urllib.error
import urllib.parse
import urllib.request
from typing import Any

from capability_probe import build_report, parse_ports


DEFAULT_MODEL = "gemini-2.5-flash"


def repo_root() -> pathlib.Path:
    return pathlib.Path(__file__).resolve().parent.parent


def gemini_ssl_context() -> ssl.SSLContext:
    """
    Build an HTTPS client context with a usable CA bundle.

    Some local Python installs on macOS do not have a working default CA path,
    even when HTTPS itself is otherwise functional. Prefer certifi when present.
    """

    try:
        import certifi  # type: ignore

        return ssl.create_default_context(cafile=certifi.where())
    except ImportError:
        return ssl.create_default_context()


def load_dotenv(path: pathlib.Path) -> None:
    if not path.exists():
        return

    for raw_line in path.read_text(encoding="utf-8").splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        key = key.strip()
        value = value.strip().strip("'").strip('"')
        os.environ.setdefault(key, value)


def load_json(path: pathlib.Path) -> dict[str, Any]:
    return json.loads(path.read_text(encoding="utf-8"))


def find_first_text_part(response: dict[str, Any]) -> str:
    candidates = response.get("candidates") or []
    for candidate in candidates:
        content = candidate.get("content") or {}
        for part in content.get("parts") or []:
            text = part.get("text")
            if isinstance(text, str) and text.strip():
                return text
    raise ValueError("Gemini response did not contain a text part")


def registration_schema() -> dict[str, Any]:
    return {
        "type": "object",
        "properties": {
            "subject_id": {"type": "string"},
            "subject_kind": {"type": "string"},
            "summary": {"type": "string"},
            "recommended_probe_strategy": {"type": "string"},
            "daemon_install_recommended": {"type": "boolean"},
            "daemon_install_reason": {"type": "string"},
            "capability_contract": {
                "type": "object",
                "properties": {
                    "subject_id": {"type": "string"},
                    "capabilities": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "name": {"type": "string"},
                                "status": {"type": "string"},
                                "commands": {
                                    "type": "array",
                                    "items": {"type": "string"},
                                },
                                "verification": {"type": "string"},
                                "evidence_summary": {"type": "string"},
                                "constraints": {
                                    "type": "array",
                                    "items": {"type": "string"},
                                },
                            },
                            "required": [
                                "name",
                                "status",
                                "commands",
                                "verification",
                                "evidence_summary",
                                "constraints",
                            ],
                        },
                    },
                },
                "required": ["subject_id", "capabilities"],
            },
            "next_actions": {
                "type": "array",
                "items": {"type": "string"},
            },
            "warnings": {
                "type": "array",
                "items": {"type": "string"},
            },
        },
        "required": [
            "subject_id",
            "subject_kind",
            "summary",
            "recommended_probe_strategy",
            "daemon_install_recommended",
            "daemon_install_reason",
            "capability_contract",
            "next_actions",
            "warnings",
        ],
    }


def build_prompt(report: dict[str, Any]) -> str:
    report_json = json.dumps(report, indent=2, sort_keys=False)
    return (
        "You are Stark, the structural authority for capability registration.\n"
        "Your job is to convert a probe report into a conservative capability registration proposal.\n\n"
        "Rules:\n"
        "- Only register capabilities supported by the probe evidence.\n"
        "- Prefer smaller true capability surfaces over inflated guesses.\n"
        "- Distinguish verified vs partial capability claims.\n"
        "- Recommend daemon install only if it clearly expands the capability surface.\n"
        "- Keep the response compact and operational.\n"
        "- Output valid JSON only.\n\n"
        "Probe report:\n"
        f"{report_json}\n"
    )


def build_protocol_prompt(report: dict[str, Any]) -> str:
    report_json = json.dumps(report, indent=2, sort_keys=False)
    return (
        "You are Stark, the structural authority for capability registration.\n"
        "Emit a small command plan using the runtime protocol shape `skyra <command_set> <command> -<args>`.\n\n"
        "Rules:\n"
        "- Output raw protocol commands only.\n"
        "- One command per line.\n"
        "- Do not output JSON.\n"
        "- Do not output markdown fences.\n"
        "- Do not explain the commands.\n"
        "- Every command must include `-reason \"...\"`.\n"
        "- Commands without `-reason` are invalid.\n"
        "- Commands must be conservative and evidence-backed.\n"
        "- Do not emit commands that assume capabilities not present in the probe report.\n"
        "- Prefer command families like `probe`, `capability`, and `registration`.\n"
        "- Use concrete args with the `-arg value` shape.\n"
        "- Keep the command list short and operational.\n\n"
        "Probe report:\n"
        f"{report_json}\n"
    )


def normalize_protocol_output(text: str) -> str:
    cleaned = text.strip()
    if cleaned.startswith("```"):
        parts = [line for line in cleaned.splitlines() if not line.strip().startswith("```")]
        cleaned = "\n".join(parts).strip()
    return cleaned


def validate_protocol_output(text: str) -> str:
    lines = [line.strip() for line in text.splitlines() if line.strip()]
    if not lines:
        raise ValueError("Protocol mode returned no commands")

    for idx, line in enumerate(lines, start=1):
        if not line.startswith("skyra "):
            raise ValueError(f"Protocol line {idx} does not start with `skyra`: {line}")
        if " -reason " not in line:
            raise ValueError(f"Protocol line {idx} is missing mandatory `-reason`: {line}")
        if not re.search(r'\s-reason\s+"[^"]+"$', line):
            raise ValueError(
                f"Protocol line {idx} must end with a quoted `-reason \"...\"` argument: {line}"
            )
    return "\n".join(lines)


def call_gemini(
    api_key: str,
    model: str,
    report: dict[str, Any],
    *,
    mode: str,
) -> Any:
    url = (
        "https://generativelanguage.googleapis.com/v1beta/models/"
        f"{urllib.parse.quote(model, safe='')}:generateContent"
    )
    if mode == "protocol":
        user_prompt = build_protocol_prompt(report)
        system_text = (
            "You emit compact runtime commands as raw text lines. "
            "Only emit commands in the form `skyra <command_set> <command> -<args>`."
        )
    else:
        user_prompt = build_prompt(report)
        response_schema = registration_schema()
        system_text = (
            "You produce evidence-backed capability registration proposals for devices. "
            "Do not invent capabilities that are not supported by the probe report."
        )

    payload = {
        "system_instruction": {
            "parts": [
                {
                    "text": system_text
                }
            ]
        },
        "contents": [
            {
                "role": "user",
                "parts": [{"text": user_prompt}],
            }
        ],
        "generationConfig": {"temperature": 0.2},
    }
    if mode != "protocol":
        payload["generationConfig"]["responseMimeType"] = "application/json"
        payload["generationConfig"]["responseJsonSchema"] = response_schema

    request = urllib.request.Request(
        url,
        data=json.dumps(payload).encode("utf-8"),
        headers={
            "Content-Type": "application/json",
            "x-goog-api-key": api_key,
        },
        method="POST",
    )

    try:
        with urllib.request.urlopen(request, timeout=45, context=gemini_ssl_context()) as response:
            raw = response.read().decode("utf-8")
    except urllib.error.HTTPError as exc:
        body = exc.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"Gemini HTTP {exc.code}: {body}") from exc
    except urllib.error.URLError as exc:
        raise RuntimeError(
            "Gemini request failed: "
            f"{exc}. If this is a local certificate issue, ensure certifi is installed "
            "or set SSL_CERT_FILE to a valid CA bundle."
        ) from exc

    parsed = json.loads(raw)
    text = find_first_text_part(parsed)
    if mode == "protocol":
        return validate_protocol_output(normalize_protocol_output(text))
    return json.loads(text)


def main() -> int:
    parser = argparse.ArgumentParser(description="Tiny Gemini-backed capability registration demo.")
    parser.add_argument(
        "--probe-report",
        help="Path to an existing probe JSON report. If omitted, the probe runs locally first.",
    )
    parser.add_argument(
        "--target",
        action="append",
        default=[],
        help="Explicit remote target to include when generating a new probe report. May be used multiple times.",
    )
    parser.add_argument(
        "--ports",
        help="Comma-separated port list for remote probing when generating a new report.",
    )
    parser.add_argument(
        "--timeout",
        type=float,
        default=0.75,
        help="Socket timeout in seconds for generated probe reports. Default: 0.75",
    )
    parser.add_argument(
        "--model",
        help=f"Gemini model override. Default: env GEMINI_MODEL or {DEFAULT_MODEL}",
    )
    parser.add_argument(
        "--output",
        help="Optional path to write the structured registration JSON.",
    )
    parser.add_argument(
        "--mode",
        choices=["registration", "protocol"],
        default="registration",
        help="registration = capability registration proposal, protocol = command-emission smoke test",
    )
    args = parser.parse_args()

    load_dotenv(repo_root() / ".env")

    api_key = os.environ.get("GEMINI_API_KEY", "").strip()
    if not api_key:
        print("error: GEMINI_API_KEY is required in the environment or .env", file=sys.stderr)
        return 2

    model = (args.model or os.environ.get("GEMINI_MODEL") or DEFAULT_MODEL).strip()
    if not model:
        print("error: model cannot be empty", file=sys.stderr)
        return 2

    if args.probe_report:
        report = load_json(pathlib.Path(args.probe_report))
    else:
        try:
            ports = parse_ports(args.ports)
        except ValueError as exc:
            print(f"error: {exc}", file=sys.stderr)
            return 2
        report = build_report(args.target, ports, args.timeout)

    registration = call_gemini(api_key, model, report, mode=args.mode)
    if args.mode == "protocol":
        payload = str(registration).strip()
        if args.output:
            pathlib.Path(args.output).write_text(payload + "\n", encoding="utf-8")
        print(payload)
    else:
        payload = json.dumps(registration, indent=2, sort_keys=False)
        if args.output:
            pathlib.Path(args.output).write_text(payload + "\n", encoding="utf-8")
        print(payload)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
