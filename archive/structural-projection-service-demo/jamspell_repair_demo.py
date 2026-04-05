#!/usr/bin/env python3
import argparse
from pathlib import Path

import jamspell

from coref_resolve import preprocess_text


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()
    parser.add_argument("text", nargs="+")
    parser.add_argument("--model", default="models/en.bin")
    return parser


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()

    text = " ".join(args.text)
    preprocessed = preprocess_text(text)

    model_path = Path(args.model)
    if not model_path.exists():
        raise SystemExit(f"Missing JamSpell model: {model_path}")

    corrector = jamspell.TSpellCorrector()
    if not corrector.LoadLangModel(str(model_path)):
        raise SystemExit(f"Failed to load JamSpell model: {model_path}")

    repaired = corrector.FixFragment(preprocessed)

    print("INPUT")
    print(text)
    print()

    if preprocessed != text:
        print("PREPROCESSED")
        print(preprocessed)
        print()

    print("REPAIRED")
    print(repaired)

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
