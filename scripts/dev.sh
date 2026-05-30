#!/bin/bash
# Script de développement avec live-reload
# Usage: ./scripts/dev.sh

if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

echo "Démarrage du serveur en mode dev avec live-reload..."
# Air va monitorer les changements dans le répertoire et redémarrer l'API
$(go env GOPATH)/bin/air -c .air.toml
