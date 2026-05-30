#!/bin/bash
# Script de migration manuel
# Usage: ./scripts/migrate.sh [up|down|version]

if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

if [ -z "$DATABASE_URL" ]; then
  echo "Erreur: DATABASE_URL n'est pas définie dans .env"
  exit 1
fi

# Utilisation de l'outil migrate (doit être installé sur la machine)
# Si non installé, on peut l'installer via : go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path migrations/ -database "$DATABASE_URL" "$@"
