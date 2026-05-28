package system

import (
	"os"
	"os/exec"
	"runtime"
)

// GetOS renvoie le nom du système d'exploitation.
func GetOS() string {
	return runtime.GOOS
}

// IsClaudeInstalled vérifie si 'claude' est dans le PATH.
func IsClaudeInstalled() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

// GetHomeDir renvoie le répertoire personnel de l'utilisateur.
func GetHomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}
