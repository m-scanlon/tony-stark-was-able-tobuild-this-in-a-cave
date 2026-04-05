#!/usr/bin/env zsh
set -euo pipefail

if [[ "$#" -eq 0 ]]; then
  echo 'Usage: ./run_dependency_demo.sh --resolve-coref "your sentence here"'
  exit 1
fi

if [[ ! -x ".venv/bin/python" ]]; then
  echo "Missing .venv. Create it and install requirements first."
  echo "  python3 -m venv .venv"
  echo "  .venv/bin/pip install -r requirements.txt"
  exit 1
fi

.venv/bin/python dependency_projection_demo.py "$@"
