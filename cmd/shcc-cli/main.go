package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/mendoc/shcc/internal/api"
	"github.com/mendoc/shcc/internal/auth"
	"github.com/mendoc/shcc/internal/config"
	"github.com/mendoc/shcc/internal/system"
)

const version = "1.0.0"

func main() {
	// 1. Initialisation des clés RSA locales
	if err := auth.EnsureKeys(); err != nil {
		fmt.Printf("Erreur d'initialisation des clés: %v\n", err)
		os.Exit(1)
	}

	// 2. Enregistrement/Mise à jour de l'utilisateur sur le serveur (Silencieux)
	go syncUser()

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
		// Logique d'update à implémenter
	case "name":
		if len(os.Args) < 3 {
			fmt.Println("Erreur: nom manquant. Usage: shcc name <nom>")
			return
		}
		handleSetName(os.Args[2])
	default:
		handleShare(cmd)
	}
}

func syncUser() {
	email, _ := config.GetUserEmail()
	dname, _ := config.GetUserDisplayName()
	pubKey, _ := auth.GetPublicKey()
	shccCfg, _ := config.ReadShccConfig()

	if email != "" && pubKey != "" {
		_ = api.RegisterUser(api.User{
			Email:       email,
			DisplayName: dname,
			Name:        shccCfg.Name,
			PublicKey:   pubKey,
		})
	}
}

func handleSetName(name string) {
	cfg, err := config.ReadShccConfig()
	if err != nil {
		fmt.Printf("Erreur lors de la lecture de la config: %v\n", err)
		return
	}
	cfg.Name = name
	if err := config.SaveShccConfig(cfg); err != nil {
		fmt.Printf("Erreur lors de la sauvegarde de la config: %v\n", err)
		return
	}

	// Forcer une synchro immédiate
	email, _ := config.GetUserEmail()
	dname, _ := config.GetUserDisplayName()
	pubKey, _ := auth.GetPublicKey()
	err = api.RegisterUser(api.User{
		Email:       email,
		DisplayName: dname,
		Name:        name,
		PublicKey:   pubKey,
	})

	if err != nil {
		fmt.Printf("Nom enregistré localement mais erreur de synchro API: %v\n", err)
	} else {
		fmt.Printf("✅ Nom '%s' enregistré et synchronisé avec succès.\n", name)
	}
}

func handleShare(recipientIdentifier string) {
	fmt.Printf("🔍 Recherche de l'utilisateur '%s'...\n", recipientIdentifier)
	
	destUser, err := api.GetUser(recipientIdentifier)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Printf("📧 Utilisateur trouvé : %s (%s)\n", destUser.DisplayName, destUser.Email)

	ownerEmail, err := config.GetUserEmail()
	if err != nil {
		fmt.Printf("❌ Erreur: impossible de récupérer votre email (%v)\n", err)
		return
	}

	creds, err := config.GetCredentials()
	if err != nil {
		fmt.Printf("❌ Erreur: impossible de lire les credentials Claude Code (%v)\n", err)
		return
	}

	fmt.Println("🔐 Chiffrement des credentials...")
	encrypted, err := auth.EncryptHybrid(destUser.PublicKey, []byte(creds))
	if err != nil {
		fmt.Printf("❌ Erreur de chiffrement : %v\n", err)
		return
	}

	encoded := base64.StdEncoding.EncodeToString(encrypted)

	// Par défaut, expiration dans 24h si non spécifié dans les credentials
	// En réalité, on devrait parser le JSON pour extraire 'expiresAt'
	payload := api.Share{
		Owner:       ownerEmail,
		To:          destUser.Email,
		Credentials: encoded,
		ExpiredAt:   time.Now().Add(24 * time.Hour),
	}

	fmt.Println("🚀 Envoi au serveur...")
	if err := api.PostShare(payload); err != nil {
		fmt.Printf("❌ Erreur lors du partage : %v\n", err)
		return
	}

	fmt.Printf("✅ Credentials partagés avec succès à %s !\n", destUser.Email)
}

func handleReceive() {
	email, err := config.GetUserEmail()
	if err != nil {
		fmt.Printf("❌ Erreur: impossible d'identifier votre compte (%v)\n", err)
		return
	}

	fmt.Println("📥 Recherche de credentials partagés pour vous...")
	shares, err := api.GetShares(email)
	if err != nil {
		fmt.Printf("❌ Erreur API : %v\n", err)
		return
	}

	if len(shares) == 0 {
		fmt.Println("📭 Aucune clé partagée trouvée.")
		return
	}

	fmt.Printf("✨ %d clé(s) trouvée(s) :\n", len(shares))
	for i, s := range shares {
		fmt.Printf("[%d] De : %s (Expire le %s)\n", i+1, s.Owner, s.ExpiredAt.Format("02/01/2006 à 15:04"))
	}

	var choice int
	fmt.Print("\nChoisissez une clé à appliquer (numéro) ou 0 pour annuler : ")
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(shares) {
		fmt.Println("Annulé.")
		return
	}

	selected := shares[choice-1]
	fmt.Println("🔓 Déchiffrement et installation...")

	decoded, err := base64.StdEncoding.DecodeString(selected.Credentials)
	if err != nil {
		fmt.Printf("❌ Erreur de décodage : %v\n", err)
		return
	}

	decrypted, err := auth.DecryptHybrid(decoded)
	if err != nil {
		fmt.Printf("❌ Échec du déchiffrement : %v (La clé n'était peut-être pas destinée à cette machine)\n", err)
		return
	}

	// Sauvegarde
	path := config.GetCredentialsPath()
	if err := os.MkdirAll(os.ExpandEnv("$HOME/.claude"), 0700); err != nil {
		fmt.Printf("❌ Erreur création dossier : %v\n", err)
		return
	}

	if err := os.WriteFile(path, decrypted, 0600); err != nil {
		fmt.Printf("❌ Erreur d'écriture : %v\n", err)
		return
	}

	fmt.Println("✅ Credentials installés avec succès ! Vous pouvez maintenant utiliser Claude Code.")
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

	dname, _ := config.GetUserDisplayName()
	if dname != "" {
		fmt.Printf("Display: %s\n", dname)
	}

	shccCfg, _ := config.ReadShccConfig()
	if shccCfg.Name != "" {
		fmt.Printf("Nom shcc:%s\n", shccCfg.Name)
	}

	fmt.Printf("Credentials Path: %s\n", config.GetCredentialsPath())
}

func printHelp() {
	fmt.Println("Utilisation: shcc <commande>")
	fmt.Println("")
	fmt.Println("Commandes:")
	fmt.Println("  status           Affiche un récap des informations détectées")
	fmt.Println("  <email> | <nom>  Partage vos credentials")
	fmt.Println("  (sans argument)  Vérifie et installe les credentials reçus")
	fmt.Println("  name <nom>       Définit un nom pour l'utilisateur courant")
	fmt.Println("  update           Met à jour le CLI shcc")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -v, --version    Affiche la version")
	fmt.Println("  -h, --help       Affiche l'aide")
}
