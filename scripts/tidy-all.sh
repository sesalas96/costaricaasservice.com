#!/usr/bin/env bash
# tidy-all.sh — go mod tidy en todos los módulos
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

modules=()
while IFS= read -r f; do
  modules+=("$(dirname "$f")")
done < <(find . -name go.mod -not -path './*/node_modules/*' -not -path './*/build/*' -print)

for mod in "${modules[@]}"; do
  echo "→ tidy $mod"
  ( cd "$mod" && go mod tidy )
done

echo "✔ tidy OK"
