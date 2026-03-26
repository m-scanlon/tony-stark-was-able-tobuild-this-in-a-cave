#!/usr/bin/env python3
import argparse
import os
import re
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import List

import torch
from transformers import AutoModelForSequenceClassification, AutoTokenizer


DEFAULT_MODEL = "typeform/distilbert-base-uncased-mnli"


@dataclass
class TripleRecord:
    sentence_index: int
    sentence_text: str
    triple_index: int
    subject: str
    relation: str
    object: str


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()
    parser.add_argument("--model", default=DEFAULT_MODEL)
    parser.add_argument("--top-k", type=int, default=0)
    return parser


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()

    demo_output = sys.stdin.read()
    triples = parse_demo_output(demo_output)
    if not triples:
        print("NLI TRIPLE SCORES")
        print("  none")
        return 0

    model_path = resolve_model_path(args.model)
    if model_path != args.model:
        force_offline_mode()

    tokenizer = AutoTokenizer.from_pretrained(model_path)
    model = AutoModelForSequenceClassification.from_pretrained(model_path)
    model.eval()

    label_lookup = {label.lower(): idx for idx, label in model.config.id2label.items()}
    entail_idx = pick_label_index(label_lookup, "entail")
    neutral_idx = pick_label_index(label_lookup, "neutral")
    contradiction_idx = pick_label_index(label_lookup, "contrad")

    scored = []
    for triple in triples:
        hypothesis = verbalize(triple)
        with torch.no_grad():
            encoded = tokenizer(
                triple.sentence_text,
                hypothesis,
                return_tensors="pt",
                truncation=True,
            )
            logits = model(**encoded).logits[0]
            probs = torch.softmax(logits, dim=-1)

        entail = float(probs[entail_idx]) if entail_idx is not None else 0.0
        neutral = float(probs[neutral_idx]) if neutral_idx is not None else 0.0
        contradiction = float(probs[contradiction_idx]) if contradiction_idx is not None else 0.0

        scored.append((triple, hypothesis, entail, neutral, contradiction))

    scored.sort(key=lambda item: item[2], reverse=True)
    if args.top_k > 0:
        scored = scored[: args.top_k]

    print("NLI TRIPLE SCORES")
    for triple, hypothesis, entail, neutral, contradiction in scored:
        print(f"  SENTENCE {triple.sentence_index} TRIPLE {triple.triple_index}")
        print(f"    premise: {triple.sentence_text}")
        print(f"    hypothesis: {hypothesis}")
        print(f"    entailment: {entail:.4f}")
        if neutral_idx is not None:
            print(f"    neutral: {neutral:.4f}")
        if contradiction_idx is not None:
            print(f"    contradiction: {contradiction:.4f}")

    return 0


def parse_demo_output(text: str) -> List[TripleRecord]:
    triples: List[TripleRecord] = []
    current_sentence_index = None
    current_sentence_text = None
    current_raw_triple_index = None
    pending_subject = None
    pending_relation = None
    pending_object = None

    for line in text.splitlines():
        sentence_match = re.match(r"^SENTENCE (\d+)$", line)
        if sentence_match:
            current_sentence_index = int(sentence_match.group(1))
            current_sentence_text = None
            current_raw_triple_index = None
            pending_subject = None
            pending_relation = None
            pending_object = None
            continue

        if current_sentence_index is not None and current_sentence_text is None and line and not line.startswith("  "):
            current_sentence_text = line.strip()
            continue

        triple_match = re.match(r"^\s{2}RAW TRIPLE (\d+)$", line)
        if triple_match:
            current_raw_triple_index = int(triple_match.group(1))
            pending_subject = None
            pending_relation = None
            pending_object = None
            continue

        if current_raw_triple_index is None:
            continue

        if line.startswith("    subject: "):
            pending_subject = line[len("    subject: ") :].strip()
        elif line.startswith("    relation: "):
            pending_relation = line[len("    relation: ") :].strip()
        elif line.startswith("    object: "):
            pending_object = line[len("    object: ") :].strip()
        elif line.startswith("    confidence: "):
            if current_sentence_text and pending_subject and pending_relation and pending_object:
                triples.append(
                    TripleRecord(
                        sentence_index=current_sentence_index,
                        sentence_text=current_sentence_text,
                        triple_index=current_raw_triple_index,
                        subject=pending_subject,
                        relation=pending_relation,
                        object=pending_object,
                    )
                )
            current_raw_triple_index = None
            pending_subject = None
            pending_relation = None
            pending_object = None

    return triples


def verbalize(triple: TripleRecord) -> str:
    return normalize_whitespace(f"{triple.subject} {triple.relation.lower()} {triple.object}.")


def normalize_whitespace(text: str) -> str:
    return " ".join(text.split())


def pick_label_index(label_lookup, prefix: str):
    for label, idx in label_lookup.items():
        if prefix in label:
            return idx
    return None


def force_offline_mode() -> None:
    os.environ.setdefault("HF_HUB_OFFLINE", "1")
    os.environ.setdefault("TRANSFORMERS_OFFLINE", "1")


def resolve_model_path(model_name: str) -> str:
    hub_root = Path.home() / ".cache" / "huggingface" / "hub"
    cache_key = "models--" + model_name.replace("/", "--")
    snapshots_dir = hub_root / cache_key / "snapshots"
    if snapshots_dir.is_dir():
        snapshots = sorted(
            [path for path in snapshots_dir.iterdir() if path.is_dir()],
            key=lambda path: path.stat().st_mtime,
            reverse=True,
        )
        for snapshot in snapshots:
            if (snapshot / "config.json").is_file():
                return str(snapshot)
    return model_name


if __name__ == "__main__":
    raise SystemExit(main())
