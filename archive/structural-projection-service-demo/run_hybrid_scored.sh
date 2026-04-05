#!/usr/bin/env zsh
set -euo pipefail

if [[ "$#" -eq 0 ]]; then
  echo 'Usage: ./run_hybrid_scored.sh --raw "your sentence here"'
  exit 1
fi

if [[ ! -x ".venv/bin/python" ]]; then
  echo "Missing .venv. Create it and install requirements first."
  echo "  python3 -m venv .venv"
  echo "  .venv/bin/pip install -r requirements.txt"
  exit 1
fi

hybrid_output=$(./run_hybrid.sh "$@")
printf '%s\n' "$hybrid_output"
printf '\n'
printf '%s\n' "$hybrid_output" | .venv/bin/python score_triples.py
