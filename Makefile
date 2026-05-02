# costaricaasservice — root Makefile
# Orquesta infra, migraciones, seed, build/test de todos los servicios.

SHELL := /usr/bin/env bash

ROOT       := $(shell pwd)
INFRA_DIR  := $(ROOT)/infra/cri-infra-docker
SCRIPTS    := $(ROOT)/scripts

# Colores para logs
C_RESET := \033[0m
C_INFO  := \033[1;34m
C_OK    := \033[1;32m
C_WARN  := \033[1;33m

.PHONY: help up down restart status logs psql dbs migrate seed build test lint tidy clean new-svc up-observability

help: ## Muestra esta ayuda
	@printf "$(C_INFO)costaricaasservice — comandos$(C_RESET)\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "} {printf "  $(C_OK)%-20s$(C_RESET) %s\n", $$1, $$2}'

up: ## Levanta infra local (postgres+kafka+redis+vault)
	@printf "$(C_INFO)→ levantando infra local...$(C_RESET)\n"
	@cd $(INFRA_DIR) && docker compose up -d

down: ## Baja infra local
	@printf "$(C_WARN)→ bajando infra local...$(C_RESET)\n"
	@cd $(INFRA_DIR) && docker compose down

restart: down up ## Reinicia infra local

status: ## Estado de los contenedores
	@cd $(INFRA_DIR) && docker compose ps

logs: ## Logs de infra (uso: make logs SVC=postgres)
	@cd $(INFRA_DIR) && docker compose logs -f $(SVC)

psql: ## Abre psql en una DB (uso: make psql DB=cri_iduc_identity)
	@docker exec -it cri-postgres psql -U postgres -d $(DB)

dbs: ## Lista las DBs creadas en postgres
	@docker exec -it cri-postgres psql -U postgres -c "\l"

migrate: ## Aplica migraciones de todos los servicios
	@bash $(SCRIPTS)/migrate-all.sh

seed: ## Carga data de prueba (1000 ciudadanos mock + members)
	@bash $(SCRIPTS)/seed-all.sh

build: ## Build de todos los servicios Go
	@bash $(SCRIPTS)/build-all.sh

test: ## go test -race en todos los módulos
	@bash $(SCRIPTS)/test-all.sh

lint: ## Lint de Go + frontends
	@bash $(SCRIPTS)/lint-all.sh

tidy: ## go mod tidy en todos los módulos
	@bash $(SCRIPTS)/tidy-all.sh

clean: ## Limpia artefactos de build y caches
	@find . -type d -name node_modules -prune -exec rm -rf {} +
	@find . -type d -name .next -prune -exec rm -rf {} +
	@find . -type d -name .build -prune -exec rm -rf {} +
	@find . -type d -name build -prune -exec rm -rf {} +
	@rm -rf .logs/ logs/

new-svc: ## Scaffold de un nuevo servicio (uso: make new-svc AREA=iduc NAME=foo PORT=8090)
	@bash $(SCRIPTS)/bootstrap-service.sh $(AREA) $(NAME) $(PORT)

up-observability: ## Levanta stack de observabilidad (loki+prometheus+grafana)
	@cd $(INFRA_DIR) && docker compose -f docker-compose.observability.yml up -d
