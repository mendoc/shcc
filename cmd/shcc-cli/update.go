package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

// GitHub Release response structure
type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func handleUpdate() {
	fmt.Println("Vérification des mises à jour...")

	// 1. Appeler l'API GitHub pour obtenir la dernière release
	resp, err := http.Get("https://api.github.com/repos/mendoc/shcc/releases/latest")
	if err != nil {
		fmt.Printf("❌ Impossible de vérifier les mises à jour : %v\n", err)
		return
	}
	defer resp.Body.Close()

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		fmt.Printf("❌ Impossible de parser la réponse GitHub : %v\n", err)
		return
	}

	latestVersion := release.TagName // Ex: v1.1.0
	currentVersion := "v" + version

	if latestVersion == currentVersion {
		fmt.Println("✅ Vous utilisez déjà la dernière version.")
		return
	}

	fmt.Printf("✨ Nouvelle version disponible : %s (Actuelle: %s)\n", latestVersion, currentVersion)
	fmt.Print("Voulez-vous mettre à jour ? (y/N) : ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		return
	}

	// 2. Relancer le script d'installation pour écraser le binaire
	cmd := exec.Command("sh", "-c", "curl -fsSL https://shcc.ongoua.pro/install.sh | bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Erreur lors de la mise à jour : %v\n", err)
	} else {
		fmt.Println("Synchronisation du profil après mise à jour...")
		syncUser()
	}
}
