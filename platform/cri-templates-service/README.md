# cri-templates-service

Plantilla de servicio Go para `costaricaasservice`. **No correr en producción** — usar como scaffold.

## Crear un servicio nuevo

```bash
./scripts/bootstrap-service.sh <area> <name> <port>
# ejemplo:
./scripts/bootstrap-service.sh iduc identity 8081
```

El script copia este directorio a `<area>/cri-svc-<name>/`, renombra módulos y reemplaza `TEMPLATE_PORT` por el puerto dado. Luego correr:

```bash
cd <area>/cri-svc-<name>
go mod tidy
go build ./...
```

## Layout

```text
cmd/
  api/          → binario HTTP del servicio
  seed/         → carga datos de prueba
config/         → {local,qa,staging,prod}.yaml
internal/
  config/       → loader Viper con override por env
  handlers/     → router chi + handlers HTTP
  service/      → lógica de negocio
  store/        → acceso a datos (pgx/v5)
  middleware/   → middlewares específicos
  event/        → producers/consumers Kafka
migrations/     → SQL up/down (golang-migrate)
tests/          → integration tests
Dockerfile
go.mod
README.md
```

## Convenciones

Ver [CLAUDE.md](../../CLAUDE.md) en raíz del repo. En particular:
- Envelope `{data, meta}` / `{error, meta}` siempre vía `httpx`.
- Multi-tenancy: realm en cada query (schema-per-realm).
- Una DB por servicio: `cri_<dominio>_<nombre>` o `cri_<nombre>`.
