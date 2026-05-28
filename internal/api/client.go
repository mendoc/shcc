package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// BaseURL est l'URL de l'API shcc. Peut être surchargée par SHCC_API_URL
var BaseURL = "https://shcc.ongoua.pro"

func init() {
	if url := os.Getenv("SHCC_API_URL"); url != "" {
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
	// On tente par email d'abord, puis par nom (l'API gère les deux via query params)
	url := fmt.Sprintf("%s/user?email=%s", BaseURL, identifier)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
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
