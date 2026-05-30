#!/bin/bash
set -e

# Usage: ./release.sh [major|minor|patch]
TYPE=${1:-minor}

# Lire la version actuelle depuis main.go
CURRENT_VERSION=$(grep 'const version' cmd/shcc-cli/main.go | cut -d'"' -f2)
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

case $TYPE in
    major) MAJOR=$((MAJOR + 1)); MINOR=0; PATCH=0 ;;
    minor) MINOR=$((MINOR + 1)); PATCH=0 ;;
    patch) PATCH=$((PATCH + 1)) ;;
    *) echo "Usage: ./release.sh [major|minor|patch]"; exit 1 ;;
esac

NEW_VERSION="$MAJOR.$MINOR.$PATCH"
echo "Bumping version $CURRENT_VERSION -> $NEW_VERSION"

# Mettre à jour main.go
sed -i "s/const version = \".*\"/const version = \"$NEW_VERSION\"/" cmd/shcc-cli/main.go

# Compiler les binaires
echo "Compilation des binaires..."
GOOS=linux GOARCH=amd64 go build -o bin/shcc-linux ./cmd/shcc-cli/main.go ./cmd/shcc-cli/update.go
GOOS=windows GOARCH=amd64 go build -o bin/shcc.exe ./cmd/shcc-cli/main.go ./cmd/shcc-cli/update.go

# Commit, Tag et Push
git add cmd/shcc-cli/main.go bin/shcc-linux bin/shcc.exe
git commit -m "Build: Bump version to $NEW_VERSION"
git tag "v$NEW_VERSION"
git push origin main
git push origin "v$NEW_VERSION"

# Créer la release GitHub
echo "Publication de la release GitHub..."
gh release create "v$NEW_VERSION" ./bin/shcc-linux ./bin/shcc.exe --title "Release $NEW_VERSION" --notes "Release automatique version $NEW_VERSION"


echo "✅ Release $NEW_VERSION publiée avec succès."
