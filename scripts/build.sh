#!/bin/bash
set -e

echo "Compilation du projet shcc..."

# Build CLI
echo "-> Compilation du CLI (cmd/shcc-cli)..."
GOOS=linux GOARCH=amd64 go build -o bin/shcc-linux ./cmd/shcc-cli/main.go ./cmd/shcc-cli/update.go

# Build API
echo "-> Compilation de l'API (cmd/shcc-api)..."
GOOS=linux GOARCH=amd64 go build -o bin/shcc-api ./cmd/shcc-api/main.go

echo "✅ Compilation terminée. Binaires disponibles dans bin/"
