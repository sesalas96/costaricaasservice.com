#!/usr/bin/env bash
# gen-dev-keys.sh — genera pares RS256 para servicios que firman JWTs en dev.
#
# NUNCA correr en prod. Las claves de prod viven en KMS (AWS / Vault).
# Este script crea PEMs en disco para que los servicios arranquen en local.
#
# Uso: ./scripts/gen-dev-keys.sh
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Servicios que firman JWTs y necesitan clave RS256.
services=(
  "iduc/cri-svc-iduc-identity"
)

if ! command -v openssl >/dev/null 2>&1; then
  echo "✘ openssl no instalado"
  exit 1
fi

for svc in "${services[@]}"; do
  out_dir="$ROOT/$svc/.keys"
  priv="$out_dir/jwt-rs256.pem"
  pub="$out_dir/jwt-rs256.pub.pem"

  if [[ -f "$priv" ]]; then
    echo "→ $svc: ya existe ($priv) — skip"
    continue
  fi

  mkdir -p "$out_dir"
  echo "→ generando RS256 keypair para $svc"

  # PKCS#8 private key, 3072 bits (mínimo recomendable hoy).
  openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:3072 -out "$priv"
  chmod 600 "$priv"
  openssl pkey -in "$priv" -pubout -out "$pub"
  echo "  ✔ $priv"
  echo "  ✔ $pub"
done

echo "✔ dev keys ready"
