#!/usr/bin/env python3
import argparse
import os
import sys
from pathlib import Path

from fastcoref import FCoref


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()
    parser.add_argument("text", nargs="+")
    parser.add_argument("--model", default="FCoref")
    parser.add_argument("--device", default="cpu")
    return parser


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()
    text = " ".join(args.text)
    preprocessed = preprocess_text(text)

    force_offline_mode()
    model = FCoref(model_name_or_path=resolve_model_path(args.model), device=args.device)
    preds = model.predict(texts=[preprocessed])
    result = preds[0]

    print("INPUT")
    print(text)
    print()

    if text != preprocessed:
        print("PREPROCESSED")
        print(preprocessed)
        print()

    print("COREFERENCE CLUSTERS")
    clusters = result.get_clusters()
    span_clusters = result.get_clusters(as_strings=False)
    if not clusters:
        print("  none")
    else:
        for i, cluster in enumerate(clusters, start=1):
            print(f"  CLUSTER {i}: {cluster}")
    print()

    resolved = resolve_text(result.text, span_clusters)
    print("RESOLVED TEXT")
    print(resolved)
    return 0


def resolve_text(text: str, clusters) -> str:
    char_clusters = []
    for cluster in clusters:
        if not cluster:
            continue
        mentions = []
        for mention in cluster:
            if isinstance(mention, (list, tuple)) and len(mention) == 2:
                start, end = mention
                mentions.append((int(start), int(end)))
        if mentions:
            char_clusters.append(mentions)

    replacements = []
    for mentions in char_clusters:
        anchor_start, anchor_end = mentions[0]
        anchor_text = text[anchor_start:anchor_end]
        for start, end in mentions[1:]:
            replacements.append((start, end, anchor_text))

    replacements.sort(key=lambda item: item[0], reverse=True)
    resolved = text
    for start, end, replacement in replacements:
      resolved = resolved[:start] + replacement + resolved[end:]
    return resolved


def preprocess_text(text: str) -> str:
    output = text
    output = output.replace("?", "? ")
    output = output.replace(".", ". ")
    output = output.replace(" u ", " you ")
    output = output.replace(" U ", " you ")
    output = output.replace("Its ", "It is ")
    output = output.replace("its ", "it is ")
    output = output.replace("Im ", "I am ")
    output = output.replace("im ", "I am ")
    output = " ".join(output.split())
    return output


def force_offline_mode() -> None:
    os.environ.setdefault("HF_HUB_OFFLINE", "1")
    os.environ.setdefault("TRANSFORMERS_OFFLINE", "1")


def resolve_model_path(model_arg: str) -> str:
    if model_arg != "FCoref":
        return model_arg

    home = Path.home()
    snapshots_dir = home / ".cache" / "huggingface" / "hub" / "models--biu-nlp--f-coref" / "snapshots"
    if snapshots_dir.is_dir():
        snapshots = sorted(
            [path for path in snapshots_dir.iterdir() if path.is_dir()],
            key=lambda path: path.stat().st_mtime,
            reverse=True,
        )
        if snapshots:
            for snapshot in snapshots:
                if (snapshot / "config.json").is_file():
                    return str(snapshot)

    return "biu-nlp/f-coref"


if __name__ == "__main__":
    raise SystemExit(main())
