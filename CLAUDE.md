# saascr — Núcleo de Conecta CR

Plataforma SaaS B2G/B2B de gobierno digital. Implementación del modelo e-Estonia (X-Road + e-ID + portal ciudadano) como producto multi-tenant.

**Stack**: Go 1.24, chi/v5, pgx/v5, Kafka (segmentio), Viper config, RS256 JWT, argon2id, Roaring bitmaps, Ed25519 (firmas ciudadanas), X.509 RSA-3072/ECDSA P-384 (members), HashiCorp Vault (KMS dev), Next.js 15, Flutter melos.
**Envelope API**: `{data, meta}` éxito / `{error, meta}` error, siempre con `requestId`.
**Ambientes**: `local`, `qa`, `staging`, `prod` (`config/{env}.yaml` + env vars vía Viper).
**Tenancy**: jerárquico — `realm` (jurisdicción soberana) → `member` (institución) → `citizen`. Schema-per-realm en Postgres.

> Para cualquier patrón de cbn ya resuelto (chi router, pgx pool, JWT signing, error envelope, Kafka producer), leer primero `/Users/devsebas/Documents/OPS/02_cbn/`. Copiar y adaptar, no rediseñar.

---

## Producto — qué se vende

| Capa | Componente saascr | Análogo Estonia |
|---|---|---|
| 1. Identidad | `iduc/` — ID ciudadana, claves, firma digital | e-ID |
| 2. Interoperabilidad | `interop/` — hub central + security server desplegable | X-Road |
| 3. Servicios institucionales | `members/` — instituciones que exponen/consumen vía interop | Registros (TSE, MIT, etc) |
| 4. Acceso ciudadano | `frontends/` — MiCR web + MiCR mobile + admin | eesti.ee + mRiik |
| Transversal | `platform/` — audit chain, notificaciones legales, files firmados | KSI Blockchain, etc |
| Entrada | `gateway/` — gateway-api + BFFs por audiencia | — |

**Garantías estonias no negociables (propuesta de valor):**
1. **Once-only** — el ciudadano da un dato al sistema una sola vez. Hacienda consulta a Registro Civil vía interop, no le pide al ciudadano.
2. **Descentralización** — cada member es dueño de sus datos. No hay BD central.
3. **Trazabilidad ciudadana** — el ciudadano ve quién consultó sus datos, cuándo y con qué propósito.
4. **Firma digital con valor jurídico** — equivalente a firma manuscrita.

---

## Estructura del Monorepo

```text
03_saascr/
├── gateway/          # API Gateway + BFFs
├── iduc/             # Identidad Digital Única (autenticación, claves, firma)
├── interop/          # Conecta — X-Road analog (hub, security-server, router, audit)
├── members/          # Instituciones (registro-civil, hacienda, ...)
├── platform/         # Audit ciudadano, notificaciones legales, files, feature flags, templates
├── frontends/        # Next.js (web) + Flutter melos (mobile)
├── infra/            # cri-infra-docker, cri-infra-k8s, cri-infra-ci-cd
├── libs/             # cri-lib-shared, cri-lib-http, cri-lib-auth, cri-lib-crypto, cri-lib-events, cri-lib-interop-client, cri-lib-observability
├── e2e/              # Playwright + integration suites
├── docs/             # ADRs, OpenAPI, runbooks, architecture
└── scripts/          # bootstrap-service.sh, migrate.sh, seed.sh, run-all.sh
```

---

## gateway/ — API Gateway y BFFs

| Repo | Descripción |
|---|---|
| `cri-gateway-api` | Gateway central. Valida JWT RS256, inyecta headers internos (`X-CRI-Sub`, `X-CRI-Roles`, `X-CRI-Realm`), rate limit por IP/user, RBAC, reverse proxy. Roaring bitmap de JTIs revocados. |
| `cri-bff-citizen` | BFF para MiCR (web + mobile). Orquesta llamadas para: ver mis datos, declaración pre-llenada, bitácora de accesos, firmar. |
| `cri-bff-member` | BFF para portales de instituciones (operadores). Catálogo de servicios, métricas de uso, invocación de servicios remotos. |
| `cri-bff-admin` | BFF para control plane saascr (gestión de realms, members, claves, billing, observabilidad). |

## iduc/ — Identidad Digital Única Costarricense

| Repo | Descripción |
|---|---|
| `cri-svc-iduc-identity` | Registro de ciudadano (cédula como identifier de realm), login email+password (argon2id) + WebAuthn/passkey, JWT RS256 (access 15min + refresh 30d rotating), revocación con jti numérico. DB: `cri_iduc_identity`. |
| `cri-svc-iduc-keys` | Generación/custodia de claves Ed25519 por ciudadano en KMS (Vault Transit en dev). Envelope encryption. Rotación, revocación. **Las claves privadas nunca salen del KMS.** DB: `cri_iduc_keys`. |
| `cri-svc-iduc-signing` | Firma digital de documentos. Validación, sellado de tiempo (RFC 3161), verificación. Recibe intent firmado del ciudadano y produce JWS/CAdES/PAdES. DB: `cri_iduc_signing`. |

## interop/ — Conecta (X-Road analog)

| Repo | Descripción |
|---|---|
| `cri-svc-interop-hub` | Registro central de members por realm, catálogo de servicios (`{member, service, version, schema_url}`), gestión de claves X.509 de members, gobernanza. DB: `cri_interop_hub`. |
| `cri-svc-security-server` | **Daemon que cada member instala/despliega**. Recibe request del svc del member, lo firma con la clave del member (JWS detached), lo enruta al security server destino. mTLS hacia el router. |
| `cri-svc-interop-router` | Plano de datos. Enruta requests entre security servers (HTTP/2 + mTLS). Idempotencia con ULID. |
| `cri-svc-interop-audit` | Bitácora hash-chained (SHA-256 + Merkle por epoch) de toda transacción inter-member. Cada request emite evento Kafka → escritura en log con hash del registro previo. Endpoint de verificación de integridad. DB: `cri_interop_audit`. |

## members/ — Instituciones

| Repo | Descripción |
|---|---|
| `cri-svc-registro-civil` | Mock TSE. Expone vía interop: `GET /persons/{cedula}`, `GET /persons/{cedula}/vital-events`. Seed de 1000 personas mock. DB: `cri_registro_civil`. |
| `cri-svc-hacienda` | Mock Hacienda. Expone vía interop: `GET /tax-status/{cedula}`, `GET /prefilled-return/{cedula}/{year}`. **Internamente consume `cri-svc-registro-civil` vía `cri-lib-interop-client` para construir la declaración pre-llenada (demuestra once-only).** DB: `cri_hacienda`. |

## platform/ — Plataforma Transversal

| Repo | Descripción |
|---|---|
| `cri-svc-audit` | Vista ciudadana del log inter-member. "Quién consultó tus datos, cuándo, con qué propósito". Lee de `cri-svc-interop-audit` y agrega para el ciudadano. DB: `cri_audit`. |
| `cri-svc-notifications` | Notificaciones electrónicas con valor legal. Entregadas a MiCR. DB: `cri_notifications`. |
| `cri-svc-files` | Archivos firmados (PDF + CAdES/PAdES). Storage local stub (S3 en prod). DB: `cri_files`. |
| `cri-svc-feature-flags` | Feature flags por realm/member. DB: `cri_feature_flags`. |
| `cri-templates-service` | **Template/scaffold** para crear nuevos servicios. Layout alineado con producción. Usar `scripts/bootstrap-service.sh`. |

## libs/ — Librerías Compartidas

| Repo | Descripción |
|---|---|
| `cri-lib-shared` | `httpx` (envelope), `errors` (AppError), `ctx` (request ID + realm), `ulid` (request_id), `pagination`. Module: `github.com/devsebas/saascr/libs/cri-lib-shared`. |
| `cri-lib-http` | Middlewares chi: logging, recover, requestId, **tenant resolver (realm)**, CORS. Cliente HTTP con interceptors. |
| `cri-lib-auth` | `auth.Principal` (Sub, Roles, Realm), JWT RS256 verifier, middlewares chi. Roles: `citizen`, `member_*`, `admin_*`. |
| `cri-lib-crypto` | Ed25519 sign/verify (firmas ciudadanas), X.509 helpers (members), SHA-256 hash chain + Merkle (audit). |
| `cri-lib-events` | Eventos Kafka (CloudEvents envelope). Producer/consumer. Topics: `cri.<domain>.events`, DLQ: `cri.dlq`. |
| `cri-lib-interop-client` | **SDK que un svc-member usa para exponer handlers y consumir servicios remotos vía security-server local.** Abstrae firma JWS, mTLS, idempotencia, retries. |
| `cri-lib-observability` | OpenTelemetry, slog (text en local, JSON en qa+), métricas Prometheus. |

---

## Convenciones

- **Organización**: `{area}/{cri-tipo-nombre}`. Áreas: `gateway`, `iduc`, `interop`, `members`, `platform`, `libs`, `infra`, `frontends`.
- **DB por servicio**: cada svc su propia DB en Postgres (`cri_<dominio>_<nombre>` o `cri_<nombre>`). **No se comparten tablas entre servicios.**
- **Schema-per-realm** (multi-tenancy): dentro de cada DB, un schema lógico por realm (e.g., `cri_iduc_identity.cr_prod.users`). NO row-level security.
- **Comunicación sync**: HTTP/JSON vía gateway para clientes externos. HTTP directo entre svcs internos del mismo realm. **Inter-member SOLO vía security-server + interop-router** (firmado y auditado).
- **Comunicación async**: Kafka eventos (CloudEvents envelope vía `cri-lib-events`). Topics: `cri.<domain>.events`. DLQ: `cri.dlq`.
- **Auth flow**: Client → Gateway (JWT verify + Roaring bitmap revocation + RBAC) → BFF → Services. Headers inyectados: `X-CRI-Sub`, `X-CRI-Roles`, `X-CRI-Realm`, `X-Request-Id`.
- **Roles** (`libs/cri-lib-auth/auth/roles.go`): `citizen`, `member_operator`, `member_admin`, `realm_admin`, `saascr_admin`.
- **Crypto**:
  - Firmas ciudadanas: Ed25519, claves en KMS (Vault Transit / AWS KMS).
  - Members: X.509 (RSA-3072 o ECDSA P-384) firmadas por el hub-CA del realm.
  - Audit: SHA-256 hash chain + Merkle root por epoch (1h o 10k entries).
  - mTLS entre security servers.
- **Idempotencia**: requests inter-member con `request_id` ULID; deduplicación en el destino por 24h.
- **Observabilidad**: OpenTelemetry desde sprint 1. Trazas propagadas vía interop. Métrica clave: p99 inter-member < 500ms.
- **Spanish-first naming** en código, comentarios y commits.

---

## Comandos rápidos

```bash
make up                           # levantar infra local (postgres+kafka+redis+vault)
make down                         # bajar infra local
make migrate                      # aplicar migraciones a todas las DBs
make seed                         # cargar 1000 ciudadanos mock + members
make test                         # go test -race + frontends test
make lint                         # golangci-lint + eslint + dart analyze
make build                        # build de todos los servicios
make psql DB=cri_iduc_identity    # abrir psql contra una DB
make new-svc AREA=iduc NAME=foo PORT=8090   # scaffoldear nuevo servicio
```

---

## Guardrails (aplican a TODOS los roles, agentes y commits)

- **NUNCA** inventar formatos de respuesta — siempre `{data, meta}` / `{error, meta}` con `requestId`.
- **NUNCA** bypass del modelo auth ni cambiar roles sin ADR.
- **NUNCA** exponer endpoints internos públicamente. **Inter-member SOLO vía security-server.**
- **NUNCA** compartir databases entre servicios.
- **NUNCA** olvidar `realm_id` en migrations o consultas. Toda migration debe nacer multi-tenant.
- **NUNCA** sacar claves privadas del KMS. La firma se hace server-side dentro del KMS.
- **SIEMPRE** propagar `X-Request-Id` end-to-end.
- **SIEMPRE** registrar acceso inter-member en `cri-svc-interop-audit` con `{requester_member, target_member, service, citizen_id?, purpose, request_hash}`.
- **SIEMPRE** verificar el catálogo (`docs/architecture/SERVICE_CATALOG.md`) antes de crear servicios nuevos.
- **SIEMPRE** validar contra `docs/architecture/SERVICE_SCAFFOLD_CHECKLIST.md` antes de declarar un servicio listo.

---

## Workflows

| Workflow | Archivo |
|---|---|
| Crear nuevo servicio | `docs/architecture/NEW_SERVICE.md` |
| Agregar member nuevo (institución) | `docs/architecture/NEW_MEMBER.md` |
| Agregar servicio inter-member (catálogo) | `docs/architecture/NEW_INTEROP_SERVICE.md` |
| Cambio de auth / identidad | `docs/architecture/AUTH_CHANGE.md` |
| Onboarding de un realm nuevo | `docs/architecture/NEW_REALM.md` |
| Respuesta a incidentes | `docs/runbooks/INCIDENT_RESPONSE.md` |

---

## Estado actual

Sprint 0 — Cimentación. MVP target: ~3 meses (6 sprints de 2 semanas). Plan vigente en [`~/.claude/plans/vamos-a-comenzar-el-ethereal-pike.md`](../../.claude/plans/vamos-a-comenzar-el-ethereal-pike.md).
