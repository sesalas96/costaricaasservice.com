# cri-infra-docker

Stack de infra local para `costaricaasservice`. Usar desde la raíz del repo: `make up`, `make down`, `make psql DB=...`.

## Servicios

| Servicio | Puerto | Notas |
|---|---|---|
| `cri-postgres` | 5432 | user/pass: `postgres`/`postgres`. Crea automáticamente las DBs (`cri_iduc_identity`, `cri_interop_hub`, …) en el primer `up` vía `postgres/init/`. |
| `cri-redis` | 6379 | Sin password en dev. |
| `cri-kafka` | 9092 | Listener PLAINTEXT_HOST para clientes en host; auto-create topics activado. |
| `cri-zookeeper` | 2181 | Coordinador de Kafka. |
| `cri-vault` | 8200 | Modo dev. Root token: `dev-root-token`. KMS para claves Ed25519 ciudadanas. |

## Crear DBs nuevas

Editar `postgres/init/01-create-databases.sh` y `docker compose down -v && make up` para recrear el volumen.

## Recursos opcionales

- `docker-compose.observability.yml` (Loki + Prometheus + Grafana + Alertmanager) — `make up-observability`.
