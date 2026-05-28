package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mendoc/shcc/internal/system"
)

const (
	keySize = 4096
)

// GetConfigDir renvoie le chemin vers ~/.shcc/
func GetConfigDir() string {
	home := system.GetHomeDir()
	return filepath.Join(home, ".shcc")
}

// EnsureKeys génère une paire de clés RSA si elles n'existent pas.
func EnsureKeys() error {
	dir := GetConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("impossible de créer le dossier de config: %w", err)
	}

	privPath := filepath.Join(dir, "id_rsa")
	if _, err := os.Stat(privPath); err == nil {
		return nil // Clés déjà présentes
	}

	// Génération
	privKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("erreur de génération RSA: %w", err)
	}

	// Sauvegarde Clé Privée
	privFile, err := os.OpenFile(privPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer privFile.Close()

	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}
	if err := pem.Encode(privFile, privBlock); err != nil {
		return err
	}

	return nil
}

// GetPublicKey renvoie la clé publique au format PEM
func GetPublicKey() (string, error) {
	privPath := filepath.Join(GetConfigDir(), "id_rsa")
	data, err := os.ReadFile(privPath)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return "", fmt.Errorf("échec du décodage PEM")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", err
	}

	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	return string(pem.EncodeToMemory(pubBlock)), nil
}

// EncryptHybrid chiffre des données en utilisant AES-GCM (clé générée) 
// et chiffre la clé AES avec la clé publique RSA.
// Format retourné : [taille_clé_rsa_chiffrée (4 octets)] + [clé_aes_chiffrée_rsa] + [nonce_aes] + [ciphertext_aes]
func EncryptHybrid(pubKeyPEM string, data []byte) ([]byte, error) {
	// 1. Charger la clé publique
	block, _ := pem.Decode([]byte(pubKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("échec du décodage de la clé publique")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("la clé n'est pas de type RSA")
	}

	// 2. Générer une clé AES 256 bits
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, err
	}

	// 3. Chiffrer la clé AES avec RSA-OAEP
	encryptedAesKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, aesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur chiffrement clé AES: %w", err)
	}

	// 4. Chiffrer les données avec AES-GCM
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// 5. Assembler le résultat
	res := append(encryptedAesKey, nonce...)
	res = append(res, ciphertext...)

	return res, nil
}

// DecryptHybrid déchiffre les données produites par EncryptHybrid
func DecryptHybrid(combined []byte) ([]byte, error) {
	privPath := filepath.Join(GetConfigDir(), "id_rsa")
	data, err := os.ReadFile(privPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("échec du décodage PEM")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 1. Extraire la clé AES chiffrée (taille fixe pour RSA 4096 = 512 octets)
	rsaKeyLen := keySize / 8
	if len(combined) < rsaKeyLen {
		return nil, fmt.Errorf("données trop courtes")
	}
	encryptedAesKey := combined[:rsaKeyLen]
	rest := combined[rsaKeyLen:]

	// 2. Déchiffrer la clé AES avec RSA
	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, encryptedAesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur déchiffrement clé AES: %w", err)
	}

	// 3. Déchiffrer les données avec AES-GCM
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(rest) < nonceSize {
		return nil, fmt.Errorf("données AES corrompues")
	}
	nonce := rest[:nonceSize]
	ciphertext := rest[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}
