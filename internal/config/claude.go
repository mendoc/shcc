package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mendoc/shcc/internal/system"
)

type ClaudeConfig struct {
	OAuthAccount struct {
		EmailAddress string `json:"emailAddress"`
		DisplayName  string `json:"displayName"`
	} `json:"oauthAccount"`
}

// ClaudeCreds définit la structure du JSON de credentials Claude Code
type ClaudeCreds struct {
	ClaudeAiOauth struct {
		ExpiresAt int64 `json:"expiresAt"`
	} `json:"claudeAiOauth"`
}

// ReadClaudeConfig lit et parse le fichier ~/.claude.json
func ReadClaudeConfig() (*ClaudeConfig, error) {
	home := system.GetHomeDir()
	configPath := filepath.Join(home, ".claude.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire %s: %w", configPath, err)
	}

	var config ClaudeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erreur de parsing %s: %w", configPath, err)
	}

	return &config, nil
}

// GetUserEmail extrait l'email de ~/.claude.json
func GetUserEmail() (string, error) {
	config, err := ReadClaudeConfig()
	if err != nil {
		return "", err
	}
	return config.OAuthAccount.EmailAddress, nil
}

// GetUserDisplayName extrait le nom d'affichage de ~/.claude.json
func GetUserDisplayName() (string, error) {
	config, err := ReadClaudeConfig()
	if err != nil {
		return "", err
	}
	return config.OAuthAccount.DisplayName, nil
}

// GetCredentialsPath renvoie le chemin vers le fichier de credentials selon l'OS
func GetCredentialsPath() string {
	home := system.GetHomeDir()
	
	// Support spécifique Windows / Linux / WSL
	// .claude est souvent un dossier caché à la racine du home
	return filepath.Join(home, ".claude", ".credentials.json")
}

// GetCredentials lit le contenu du fichier de credentials
func GetCredentials() (string, error) {
	if runtime.GOOS == "darwin" {
		return "", fmt.Errorf("l'accès automatique au Keychain macOS n'est pas implémenté. Veuillez copier vos credentials dans %s", GetCredentialsPath())
	}

	path := GetCredentialsPath()
	
	// Vérification explicite de l'existence
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("fichier de credentials introuvable à : %s. Veuillez vous assurer que Claude Code est installé et configuré", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("impossible de lire les credentials à %s: %w", path, err)
	}

	return string(data), nil
}
