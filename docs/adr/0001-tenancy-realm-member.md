# ADR-0001 — Tenancy: Realm + Member + Citizen, schema-per-realm

- **Estado:** Aceptado
- **Fecha:** 2026-05-02
- **Decisores:** equipo costaricaasservice

## Contexto

`costaricaasservice` se vende a gobiernos (cada uno como una jurisdicción soberana) y a redes de instituciones. El producto debe soportar varios "países" o "estados" en la misma instalación SaaS sin que un breach o auditoría en uno toque a otro. Además debe soportar federaciones tipo Estonia–Finlandia (NIIS) en una fase posterior.

## Decisión

Modelo de tenancy en tres niveles, jerárquico:

1. **Realm** = jurisdicción soberana (`cr-prod`, `sv-prod`, `demo`).
   - **Aislamiento fuerte.** Un realm = un schema lógico por DB.
   - Claves criptográficas por realm (CA del hub-CA por realm).
   - Cada realm puede tener su propio FQDN (`cr.costaricaasservice.io`, `sv.costaricaasservice.io`).
2. **Member** = institución dentro de un realm (Registro Civil, Hacienda, Banco Nacional).
   - Cada member tiene su par de claves X.509 firmadas por la CA del realm.
   - Cada member opera (idealmente despliega) su propio `cri-svc-security-server`.
3. **Citizen** = persona física registrada en un realm.
   - Identificada por la cédula/DNI dentro del realm. **Una cédula no es global**, su unicidad es por realm.
   - Claves Ed25519 por ciudadano custodiadas en KMS, scoped al realm.

### Implementación: schema-per-realm

- En cada DB de servicio (`cri_iduc_identity`, `cri_interop_hub`, etc.), crear un schema por realm:
  - `cri_iduc_identity.cr_prod`, `cri_iduc_identity.sv_prod`, `cri_iduc_identity.demo`.
- Las migraciones se aplican por schema. El servicio recibe `realm` por header (`X-CRI-Realm`) o claim del JWT, y usa `SET search_path TO <realm>;` al iniciar la conexión.
- **NO row-level security** — históricamente frágil y difícil de auditar.
- **NO una DB por realm** — explosión operacional cuando el costaricaasservice opere 5+ realms.

### Resolución del realm en cada request

Orden de resolución (primer hit gana):

1. JWT claim `realm` (preferido — autoritativo).
2. Header `X-CRI-Realm` (válido solo en endpoints internos behind gateway).
3. Subdominio (`{realm}.costaricaasservice.io`) — solo público.

Middleware en `cri-lib-http` lo resuelve y lo pone en `ctx`.

### Inter-realm

**Vetado en MVP.** Los servicios de un realm no pueden consumir datos de otro realm. La federación cross-realm es una decisión de Fase C y requiere su propio ADR.

## Consecuencias

**Positivas**
- Aislamiento real entre clientes/jurisdicciones.
- Auditorías por realm son simples (schema + logs filtrados).
- Permite escalar a federación cross-realm más adelante (con su propio puente firmado).

**Negativas**
- Migraciones más complejas — hay que aplicarlas por schema.
- El driver pgx debe `SET search_path` al adquirir conexión del pool (custom acquire hook).
- Tooling (psql) requiere `\c <db>; SET search_path TO <realm>;` para inspeccionar manualmente.

## Alternativas descartadas

- **DB-per-realm**: explosión operacional, hace 50+ DBs cuando hay 10 realms × 5 servicios.
- **Row-level security**: difícil de auditar, fácil olvidar el filtro, performance impredecible.
- **Tenant column en cada tabla**: nada más que RLS sin las protecciones del motor.

## Referencias

- Estonia / X-Road federation pattern (NIIS).
- Postgres `search_path` semantics.
