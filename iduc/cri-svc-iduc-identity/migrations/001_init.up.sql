-- Migration 001 — bootstrap del schema de IDUC identity para un realm.
-- Este archivo se aplica con `search_path` ya configurado al schema del realm
-- (ej. `SET search_path TO cr_prod;`). El runner de migraciones es
-- responsable de crear el schema antes de aplicar (`CREATE SCHEMA IF NOT EXISTS`).

-- Ciudadanos del realm. La cédula es única dentro del realm.
CREATE TABLE IF NOT EXISTS citizens (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cedula          TEXT NOT NULL,
  email           TEXT NOT NULL,
  password_hash   TEXT NOT NULL,                      -- argon2id encoded string
  status          TEXT NOT NULL DEFAULT 'active',     -- active | suspended | deleted
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT citizens_cedula_unique UNIQUE (cedula),
  CONSTRAINT citizens_email_unique  UNIQUE (email)
);

-- Sequence para JTI numéricos. Permite al gateway usar Roaring bitmap
-- de IDs revocados (ver ADR-0001 / cbn).
CREATE SEQUENCE IF NOT EXISTS jti_seq AS BIGINT INCREMENT BY 1 START WITH 1;

-- Sesiones activas (refresh tokens rotating). El refresh_token_hash es el
-- SHA-256 del refresh token entregado al cliente; nunca guardamos el token
-- en claro.
CREATE TABLE IF NOT EXISTS sessions (
  id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  citizen_id         UUID NOT NULL REFERENCES citizens(id) ON DELETE CASCADE,
  refresh_token_hash BYTEA NOT NULL,
  jti                BIGINT NOT NULL,                 -- JTI del access token actual emitido contra esta sesión
  device             TEXT,
  ip                 INET,
  issued_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at         TIMESTAMPTZ NOT NULL,
  revoked_at         TIMESTAMPTZ,
  rotated_to         UUID REFERENCES sessions(id),    -- detección de reuse de refresh

  CONSTRAINT sessions_refresh_unique UNIQUE (refresh_token_hash)
);

CREATE INDEX IF NOT EXISTS sessions_citizen_active_idx
  ON sessions (citizen_id) WHERE revoked_at IS NULL;

-- Lista de JTI revocados. El gateway hace polling de
-- /internal/revoked-jti/snapshot y reconstruye su Roaring bitmap.
CREATE TABLE IF NOT EXISTS revoked_jtis (
  jti        BIGINT PRIMARY KEY,
  revoked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL                     -- después de esto el JTI puede limpiarse: el access ya caducó
);

CREATE INDEX IF NOT EXISTS revoked_jtis_expires_idx ON revoked_jtis (expires_at);
