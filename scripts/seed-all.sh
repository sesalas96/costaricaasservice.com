#!/usr/bin/env bash
# seed-all.sh — placeholder: invoca cmd/seed/ de cada servicio que lo defina.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

seeds=()
while IFS= read -r d; do
  seeds+=("$d")
done < <(find iduc interop members platform -path '*/cmd/seed' -type d 2>/dev/null || true)

if [[ ${#seeds[@]} -eq 0 ]]; then
  echo "→ no seed commands yet"
  exit 0
fi

for seed_dir in "${seeds[@]}"; do
  svc_dir="$(dirname "$(dirname "$seed_dir")")"
  echo "→ seeding $svc_dir"
  ( cd "$svc_dir" && go run ./cmd/seed )
done

echo "✔ seed OK"
