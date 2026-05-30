#!/bin/bash
set -e

echo "Compilation du projet shcc..."

# Build CLI Linux
echo "-> Compilation du CLI Linux (bin/shcc-linux)..."
GOOS=linux GOARCH=amd64 go build -o bin/shcc-linux ./cmd/shcc-cli/main.go ./cmd/shcc-cli/update.go

# Build CLI Windows
echo "-> Compilation du CLI Windows (bin/shcc.exe)..."
GOOS=windows GOARCH=amd64 go build -o bin/shcc.exe ./cmd/shcc-cli/main.go ./cmd/shcc-cli/update.go

# Build API
echo "-> Compilation de l'API (bin/shcc-api)..."
GOOS=linux GOARCH=amd64 go build -o bin/shcc-api ./cmd/shcc-api/main.go

echo "✅ Compilation terminée. Binaires disponibles dans bin/"
