package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/shcc/shcc/internal/system"
)

type ClaudeConfig struct {
	OAuthAccount struct {
		EmailAddress string `json:"emailAddress"`
	} `json:"oauthAccount"`
}

// GetUserEmail extrait l'email de ~/.claude.json
func GetUserEmail() (string, error) {
	home := system.GetHomeDir()
	configPath := filepath.Join(home, ".claude.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("impossible de lire %s: %w", configPath, err)
	}

	var config ClaudeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("erreur de parsing %s: %w", configPath, err)
	}

	return config.OAuthAccount.EmailAddress, nil
}

// GetCredentials lit le contenu du fichier de credentials
func GetCredentials() (string, error) {
	path := GetCredentialsPath()
	if runtime.GOOS == "darwin" && path == "macOS Keychain (Claude Code-credentials)" {
		// Logique Keychain à implémenter plus tard si nécessaire
		return "", fmt.Errorf("lecture Keychain non implémentée")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("impossible de lire les credentials à %s: %w", path, err)
	}

	return string(data), nil
}

// GetCredentialsPath renvoie le chemin vers le fichier de credentials selon l'OS
func GetCredentialsPath() string {
	home := system.GetHomeDir()
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, ".claude", ".credentials.json")
	case "darwin":
		return "macOS Keychain (Claude Code-credentials)"
	default: // linux, wsl
		return filepath.Join(home, ".claude", ".credentials.json")
	}
}
