package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/shcc/shcc/internal/api"
	"github.com/shcc/shcc/internal/auth"
	"github.com/shcc/shcc/internal/config"
	"github.com/shcc/shcc/internal/system"
)

const version = "1.0.0"

func main() {
	// Initialisation des clés si nécessaire
	if err := auth.EnsureKeys(); err != nil {
		fmt.Printf("Erreur d'initialisation des clés: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		handleReceive()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "-v", "--version":
		fmt.Printf("shcc version %s\n", version)
	case "status":
		handleStatus()
	case "-h", "--help":
		printHelp()
	case "update":
		fmt.Println("Vérification des mises à jour...")
		// Logique d'update
	case "name":
		if len(os.Args) < 3 {
			fmt.Println("Erreur: nom manquant. Usage: shcc name <nom>")
			return
		}
		fmt.Printf("Nom défini : %s (Simulé)\n", os.Args[2])
	default:
		// Si l'argument ressemble à un email ou un nom, on partage
		handleShare(cmd)
	}
}

func handleShare(recipient string) {
	fmt.Printf("Préparation du partage pour : %s\n", recipient)

	// 1. Récupérer l'email de l'owner
	owner, err := config.GetUserEmail()
	if err != nil {
		fmt.Printf("Erreur: impossible de récupérer votre email (%v)\n", err)
		return
	}

	// 2. Récupérer les credentials
	creds, err := config.GetCredentials()
	if err != nil {
		fmt.Printf("Erreur: impossible de lire les credentials Claude Code (%v)\n", err)
		return
	}

	// 3. Récupérer la clé publique du destinataire (MOCK)
	// Pour le test, on utilise notre propre clé publique pour simuler un "auto-partage"
	pubKey, err := auth.GetPublicKey()
	if err != nil {
		fmt.Printf("Erreur: impossible d'obtenir la clé publique (%v)\n", err)
		return
	}

	// 4. Chiffrer (Hybride AES + RSA)
	encrypted, err := auth.EncryptHybrid(pubKey, []byte(creds))
	if err != nil {
		fmt.Printf("Erreur de chiffrement : %v\n", err)
		return
	}

	// 5. Encoder en Base64
	encoded := base64.StdEncoding.EncodeToString(encrypted)

	// 6. Envoyer (MOCK)
	payload := api.SharePayload{
		Owner:       owner,
		To:          recipient,
		Credentials: encoded,
	}

	if err := api.MockPostShare(payload); err != nil {
		fmt.Printf("Erreur lors de l'envoi : %v\n", err)
		return
	}

	fmt.Println("✅ Credentials partagés avec succès !")
}

func handleReceive() {
	fmt.Println("Recherche de credentials partagés pour vous...")
	// Logique interactive simulée
	fmt.Println("Aucune clé partagée trouvée en base de données (Simulé).")
}

func handleStatus() {
	fmt.Println("--- shcc Status ---")
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("OS:      %s\n", system.GetOS())

	if system.IsClaudeInstalled() {
		fmt.Println("Claude Code: Installé")
	} else {
		fmt.Println("Claude Code: Non détecté")
		fmt.Println("👉 Consultez la documentation pour l'installer: https://docs.anthropic.com/claude/docs/claude-code")
	}

	email, err := config.GetUserEmail()
	if err != nil {
		fmt.Printf("Email:   Non trouvé (%v)\n", err)
	} else {
		fmt.Printf("Email:   %s\n", email)
	}

	fmt.Printf("Credentials Path: %s\n", config.GetCredentialsPath())
}

func printHelp() {
	fmt.Println("Utilisation: shcc <commande>")
	fmt.Println("")
	fmt.Println("Commandes:")
	fmt.Println("  status           Affiche un récap des informations détectées")
	fmt.Println("  <email> | <nom>  Partage vos credentials")
	fmt.Println("  update           Met à jour le CLI shcc")
	fmt.Println("  name <nom>       Définit un nom pour l'utilisateur courant")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -v, --version    Affiche la version")
	fmt.Println("  -h, --help       Affiche l'aide")
}
