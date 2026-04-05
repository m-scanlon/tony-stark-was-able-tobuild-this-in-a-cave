#!/usr/bin/env zsh
set -euo pipefail

if [[ "$#" -eq 0 ]]; then
  echo 'Usage: ./run.sh "your sentence here"'
  exit 1
fi

exec_args=$(python3 -c 'import shlex,sys; print(" ".join(shlex.quote(a) for a in sys.argv[1:]))' "$@")

mvn -q compile exec:java \
  -Dexec.mainClass=StructuralProjectionDemo \
  -Dexec.args="$exec_args"
