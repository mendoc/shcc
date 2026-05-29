package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mendoc/shcc/internal/system"
)

// BaseURL est l'URL de l'API shcc. Doit être définie via SHCC_API_URL
var BaseURL = "https://shcc.ongoua.pro"

func init() {
	if url := os.Getenv("SHCC_API_URL"); url != "" {
		system.Debug("API client init, override BaseURL to %s", url)
		BaseURL = url
	}
}

// RegisterUser enregistre ou met à jour l'utilisateur courant sur le serveur
func RegisterUser(user User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := http.Post(BaseURL+"/user", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("impossible de joindre le serveur API (%s): %w", BaseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("erreur API: le serveur est introuvable (404)")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erreur serveur (code %d) lors de l'enregistrement", resp.StatusCode)
	}

	return nil
}

// GetUser récupère les infos d'un utilisateur par email ou nom
func GetUser(identifier string) (*User, error) {
	url := fmt.Sprintf("%s/user?email=%s", BaseURL, identifier)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Tentative par nom
		url = fmt.Sprintf("%s/user?name=%s", BaseURL, identifier)
		resp, err = http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("utilisateur '%s' introuvable", identifier)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur API (code %d)", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		system.Debug("%s", err)
		return nil, fmt.Errorf("réponse serveur invalide (format JSON attendu)")
	}

	return &user, nil
}

// PostShare envoie les credentials chiffrés au serveur
func PostShare(payload Share) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(BaseURL+"/share", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("erreur API: le serveur est introuvable (404)")
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erreur API lors du partage (code %d)", resp.StatusCode)
	}

	return nil
}

// GetShares récupère les partages pour l'utilisateur courant
func GetShares(email string) ([]Share, error) {
	url := fmt.Sprintf("%s/share?to=%s", BaseURL, email)
	system.Debug("Appeler %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("erreur API: le serveur est introuvable (404)")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur serveur lors de la récupération (code %d)", resp.StatusCode)
	}

	var shares []Share
	if err := json.NewDecoder(resp.Body).Decode(&shares); err != nil {
		system.Debug("%s", err)
		return nil, fmt.Errorf("réponse serveur invalide (format JSON attendu)")
	}

	return shares, nil
}
