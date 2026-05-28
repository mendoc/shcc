package api

import (
	"fmt"
)

// SharePayload représente les données envoyées à /share
type SharePayload struct {
	Owner       string `json:"owner"`
	To          string `json:"to"`
	Credentials string `json:"credentials"` // Base64 du JSON chiffré
}

// MockGetPublicKey simule la récupération d'une clé publique pour un destinataire
func MockGetPublicKey(recipient string) (string, error) {
	// En situation réelle, on ferait un GET /user?email=recipient
	// Pour le test, on va renvoyer notre propre clé publique (auto-partage simulant)
	return "", fmt.Errorf("API non disponible (mock)")
}

// MockPostShare simule l'appel à POST /share
func MockPostShare(payload SharePayload) error {
	fmt.Printf("Simulant POST /share :\n")
	fmt.Printf("  De : %s\n", payload.Owner)
	fmt.Printf("  À : %s\n", payload.To)
	fmt.Printf("  Data (Base64) : %s...\n", payload.Credentials[:20])
	return nil
}
