#!/usr/bin/env bash
# build-all.sh — go build de todos los módulos del monorepo
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

modules=()
while IFS= read -r f; do
  modules+=("$(dirname "$f")")
done < <(find . -name go.mod -not -path './*/node_modules/*' -not -path './*/build/*' -print)

if [[ ${#modules[@]} -eq 0 ]]; then
  echo "→ no go modules found yet"
  exit 0
fi

for mod in "${modules[@]}"; do
  echo "→ building $mod"
  ( cd "$mod" && go build ./... )
done

echo "✔ build OK"
