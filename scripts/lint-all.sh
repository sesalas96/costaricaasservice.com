#!/usr/bin/env bash
# lint-all.sh — golangci-lint en todos los módulos Go; eslint+dart en frontends si existen
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

modules=()
while IFS= read -r f; do
  modules+=("$(dirname "$f")")
done < <(find . -name go.mod -not -path './*/node_modules/*' -not -path './*/build/*' -print)

if command -v golangci-lint >/dev/null 2>&1; then
  for mod in "${modules[@]}"; do
    echo "→ linting $mod"
    ( cd "$mod" && golangci-lint run ./... )
  done
else
  echo "⚠ golangci-lint no instalado — skipping go lint"
fi

# Frontends — opcional, solo si existen
if [[ -d frontends/cri-web-micr ]]; then
  ( cd frontends/cri-web-micr && npm run lint --if-present )
fi

echo "✔ lint OK"
