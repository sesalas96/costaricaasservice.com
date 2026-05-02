# Migrations — cri-svc-iduc-identity

Las migraciones se aplican **por realm**. El runner es responsable de:

1. `CREATE SCHEMA IF NOT EXISTS <realm_schema>;` (e.g., `cr_prod`).
2. `SET search_path TO <realm_schema>;`.
3. Aplicar cada `NNN_*.up.sql` en orden.

Para dev local, `make migrate` ejecuta el runner (`scripts/migrate-all.sh`) que itera por los realms configurados en `infra/cri-infra-docker/envs/realms.txt` y aplica cada migration usando `golang-migrate` con `-path migrations/`.

> **Importante:** las migrations NO crean el schema ni hacen `SET search_path`; eso es responsabilidad del runner. Esto las hace reutilizables entre realms y evita que un drop accidental impacte a otro tenant.
