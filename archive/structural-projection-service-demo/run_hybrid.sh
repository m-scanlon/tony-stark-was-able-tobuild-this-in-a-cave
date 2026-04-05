#!/usr/bin/env zsh
set -euo pipefail

if [[ "$#" -eq 0 ]]; then
  echo 'Usage: ./run_hybrid.sh --theme software_work --theme-note "..." "your sentence here"'
  exit 1
fi

theme_args=()
text_args=()

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --theme|--theme-note)
      if [[ "$#" -lt 2 ]]; then
        echo "Missing value for $1"
        exit 1
      fi
      theme_args+=("$1" "$2")
      shift 2
      ;;
    --raw)
      theme_args+=("$1")
      shift
      ;;
    *)
      text_args+=("$1")
      shift
      ;;
  esac
done

if [[ "${#text_args[@]}" -eq 0 ]]; then
  echo 'Usage: ./run_hybrid.sh --theme software_work --theme-note "..." "your sentence here"'
  exit 1
fi

if [[ ! -x ".venv/bin/python" ]]; then
  echo "Missing .venv. Create it and install requirements first."
  echo "  python3 -m venv .venv"
  echo "  .venv/bin/pip install -r requirements.txt"
  exit 1
fi

resolved_output=$(.venv/bin/python coref_resolve.py "${text_args[@]}")
printf '%s\n\n' "$resolved_output"

resolved_text=$(printf '%s\n' "$resolved_output" | awk 'found {print} /^RESOLVED TEXT$/ {found=1; next}')
if [[ -z "${resolved_text// }" ]]; then
  echo "Failed to extract resolved text from coref output."
  exit 1
fi

./run.sh "${theme_args[@]}" "$resolved_text"
