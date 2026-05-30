#!/bin/bash
# Script de migration manuel
# Usage: ./scripts/migrate.sh [up|down|version]

# Chargement du .env depuis la racine du projet
ENV_PATH="../.env"

if [ -f "$ENV_PATH" ]; then
  # On lit le fichier et on exporte les variables (en ignorant les commentaires et lignes vides)
  export $(grep -v '^#' "$ENV_PATH" | xargs)
else
  echo "Erreur: Fichier .env introuvable à $ENV_PATH"
  exit 1
fi

if [ -z "$DATABASE_URL" ]; then
  echo "Erreur: DATABASE_URL n'est pas définie dans .env"
  exit 1
fi

# Utilisation de l'outil migrate
migrate -path ../migrations/ -database "$DATABASE_URL" "$@"
