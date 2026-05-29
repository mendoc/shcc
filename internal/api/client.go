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
var BaseURL string

func init() {
	url := os.Getenv("SHCC_API_URL")
	if url == "" {
		// Pas de panic ici car init() est appelé très tôt, 
		// on laisse le CLI ou le serveur échouer lors de l'usage.
		system.Debug("SHCC_API_URL non définie")
	} else {
		system.Debug("API client init, setting BaseURL to %s", url)
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
		return fmt.Errorf("erreur de connexion à l'API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erreur API (%d)", resp.StatusCode)
	}

	return nil
}

// GetUser récupère les infos d'un utilisateur par email ou nom
func GetUser(identifier string) (*User, error) {
	url := fmt.Sprintf("%s/user?email=%s", BaseURL, identifier)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		url = fmt.Sprintf("%s/user?name=%s", BaseURL, identifier)
		resp, err = http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("utilisateur '%s' introuvable", identifier)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
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

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erreur API lors du partage (%d)", resp.StatusCode)
	}

	return nil
}

// GetShares récupère les partages pour l'utilisateur courant
func GetShares(email string) ([]Share, error) {
	url := fmt.Sprintf("%s/share?to=%s", BaseURL, email)
	system.Debug("Appeler %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur API lors de la récupération (%d)", resp.StatusCode)
	}

	var shares []Share
	if err := json.NewDecoder(resp.Body).Decode(&shares); err != nil {
		return nil, err
	}

	return shares, nil
}
