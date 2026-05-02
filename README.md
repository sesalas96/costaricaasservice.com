# costaricaasservice — Núcleo de Conecta CR

Plataforma SaaS B2G/B2B de gobierno digital, inspirada en el modelo e-Estonia.

Tres capas de producto:

1. **Conecta** — interoperabilidad inter-institucional (X-Road analog).
2. **IDUC** — Identidad Digital Única (e-ID analog).
3. **MiCR** — Carpeta Digital Ciudadana (eesti.ee + mRiik analog).

## Stack

Go 1.24 (chi/v5, pgx/v5, Kafka via segmentio, RS256 JWT, argon2id) — Next.js 15 (App Router) — Flutter (melos workspace) — PostgreSQL — Kafka — Redis — HashiCorp Vault (dev) / AWS KMS (prod).

## Dónde empezar

- [`CLAUDE.md`](CLAUDE.md) — convenciones del monorepo y catálogo de servicios.
- [`docs/adr/`](docs/adr/) — decisiones arquitectónicas (tenancy, audit, KMS).
- [`infra/cri-infra-docker/`](infra/cri-infra-docker/) — `make up` levanta el stack local.
- [`platform/cri-templates-service/`](platform/cri-templates-service/) — esqueleto de referencia para nuevos servicios.

## Comandos rápidos

```bash
make up          # levantar infra local (postgres, kafka, redis, vault)
make down        # bajar infra local
make migrate     # aplicar migraciones a todas las DBs
make seed        # cargar data de prueba (1000 ciudadanos mock)
make test        # correr tests Go con -race
make lint        # lint Go + frontends
make build       # build de todos los servicios
make psql DB=cri_iduc_identity   # abrir psql en una DB
```

## Estado

Sprint 0 — Cimentación. Ver [el plan vigente](../.claude/plans/vamos-a-comenzar-el-ethereal-pike.md).
# costaricaasservice.com
