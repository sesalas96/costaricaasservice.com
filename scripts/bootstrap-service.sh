#!/usr/bin/env bash
# bootstrap-service.sh — scaffoldea un nuevo svc copiando platform/cri-templates-service.
# Uso: ./scripts/bootstrap-service.sh <area> <name> <port>
# Ejemplo: ./scripts/bootstrap-service.sh iduc identity 8081
set -euo pipefail

if [[ $# -ne 3 ]]; then
  echo "uso: $0 <area> <name> <port>"
  echo "areas: iduc, interop, members, platform, gateway"
  exit 2
fi

AREA="$1"
NAME="$2"
PORT="$3"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEMPLATE="$ROOT/platform/cri-templates-service"
TARGET="$ROOT/$AREA/cri-svc-$NAME"

if [[ ! -d "$TEMPLATE" ]]; then
  echo "✘ template no existe en $TEMPLATE"
  exit 1
fi

if [[ -e "$TARGET" ]]; then
  echo "✘ target ya existe: $TARGET"
  exit 1
fi

echo "→ scaffolding $TARGET desde $TEMPLATE"
cp -R "$TEMPLATE" "$TARGET"

# Renombrar:
#   - el path de área:           platform/cri-templates-service → <area>/cri-svc-<name>
#   - el nombre del binario/svc: cri-templates-service          → cri-svc-<name>
#   - el prefijo de DB:          cri_templates                  → cri_<name-con-underscores>
#   - el placeholder de puerto:  TEMPLATE_PORT                  → <port>
NAME_DB="${NAME//-/_}"
find "$TARGET" -type f \( -name "*.go" -o -name "go.mod" -o -name "*.yaml" -o -name "Dockerfile" -o -name "README.md" \) -exec \
  sed -i.bak \
    -e "s|saascr/platform/cri-templates-service|saascr/$AREA/cri-svc-$NAME|g" \
    -e "s|cri-templates-service|cri-svc-$NAME|g" \
    -e "s|cri_templates|cri_$NAME_DB|g" \
    -e "s|TEMPLATE_PORT|$PORT|g" \
    {} \;
find "$TARGET" -name "*.bak" -delete

echo "✔ creado $TARGET (puerto $PORT)"
echo "  siguiente: cd $TARGET && go mod tidy && go build ./..."
