-- Migration 001 — bootstrap schema del cri-svc-interop-hub para un realm.
-- El runner crea CREATE SCHEMA IF NOT EXISTS <realm> y SET search_path antes
-- de aplicar este archivo (mismo patrón que iduc-identity).

-- Members del realm (instituciones que participan en el catálogo).
CREATE TABLE IF NOT EXISTS members (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  member_id    TEXT NOT NULL,                       -- slug estable, único por realm (ej: "registro-civil")
  display_name TEXT NOT NULL,
  description  TEXT NOT NULL DEFAULT '',
  public_key   TEXT NOT NULL,                       -- PEM PKIX (RSA-3072 o ECDSA P-384)
  status       TEXT NOT NULL DEFAULT 'active',      -- active | suspended | revoked
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT members_member_id_unique UNIQUE (member_id)
);

CREATE INDEX IF NOT EXISTS members_status_idx ON members (status);

-- Servicios publicados al catálogo. Cada servicio pertenece a un member y
-- tiene una versión semver. El (member_id, service_id, version) es único.
CREATE TABLE IF NOT EXISTS services (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  member_id   UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
  service_id  TEXT NOT NULL,                        -- slug ("persons.get", "tax.prefill")
  version     TEXT NOT NULL DEFAULT 'v1',
  description TEXT NOT NULL DEFAULT '',
  schema_url  TEXT NOT NULL DEFAULT '',             -- OpenAPI/JSON Schema URL externa
  exposed     BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT services_unique UNIQUE (member_id, service_id, version)
);

CREATE INDEX IF NOT EXISTS services_member_idx ON services (member_id) WHERE exposed = TRUE;
