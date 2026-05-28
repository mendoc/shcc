package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Récupération du port via variable d'env (standard Cloud Run)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/share", handleShare)
	http.HandleFunc("/user", handleUser)

	log.Printf("shcc-api à l'écoute sur le port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "shcc-api",
	})
}

func handleShare(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Fprintln(w, "Endpoint POST /share (à implémenter)")
		return
	}
	if r.Method == http.MethodGet {
		fmt.Fprintln(w, "Endpoint GET /share (à implémenter)")
		return
	}
	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	// Pour récupérer la clé publique de l'utilisateur
	fmt.Fprintln(w, "Endpoint /user (à implémenter)")
}
