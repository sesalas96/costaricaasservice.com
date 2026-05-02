#!/usr/bin/env bash
# migrate-all.sh — placeholder: aplica migraciones por servicio.
# Se reescribirá cuando lleguen las primeras migraciones reales.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

services=()
while IFS= read -r d; do
  services+=("$d")
done < <(find iduc interop members platform gateway -mindepth 1 -maxdepth 1 -type d 2>/dev/null || true)

if [[ ${#services[@]} -eq 0 ]]; then
  echo "→ no services with migrations yet"
  exit 0
fi

for svc in "${services[@]}"; do
  if [[ -d "$svc/migrations" ]]; then
    echo "→ TODO: applying migrations of $svc (golang-migrate per realm)"
  fi
done

echo "⚠ migrate-all.sh es un stub — implementar con golang-migrate cuando haya migraciones reales"
