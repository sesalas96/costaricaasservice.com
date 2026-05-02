#!/usr/bin/env bash
# Crea las DBs por servicio. Cada servicio tiene su propia DB; dentro de cada DB
# se aplican schemas por realm (multi-tenancy schema-per-realm).
set -euo pipefail

DBS=(
  cri_iduc_identity
  cri_iduc_keys
  cri_iduc_signing
  cri_interop_hub
  cri_interop_audit
  cri_registro_civil
  cri_hacienda
  cri_audit
  cri_notifications
  cri_files
  cri_feature_flags
)

for db in "${DBS[@]}"; do
  echo "→ creating database: $db"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE $db;
EOSQL
done

echo "✔ all databases created"
