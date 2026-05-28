#!/bin/bash
set -e

REPO_URL="https://shcc.ongoua.pro/bin/linux/shcc"
INSTALL_DIR="/usr/local/bin"

echo "Installation de shcc..."
if [ -w "$INSTALL_DIR" ]; then
    curl -fsSL "$REPO_URL" -o "$INSTALL_DIR/shcc"
    chmod +x "$INSTALL_DIR/shcc"
else
    sudo curl -fsSL "$REPO_URL" -o "$INSTALL_DIR/shcc"
    sudo chmod +x "$INSTALL_DIR/shcc"
fi

echo "✅ shcc a été installé avec succès dans $INSTALL_DIR"
shcc status
