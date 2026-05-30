package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/mendoc/shcc/internal/api"
	"github.com/mendoc/shcc/internal/database"
)

func main() {
	// Connexion à la base de données
	if err := database.Connect(); err != nil {
		log.Printf("Attention: Connexion DB échouée: %v (L'API tournera en mode dégradé)", err)
	} else {
		defer database.Close()
		log.Println("Connecté à PostgreSQL")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/share", handleShare)
	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/install.sh", handleInstallScript)
	http.Handle("/bin/", http.StripPrefix("/bin/", http.FileServer(http.Dir("./bin"))))

	log.Printf("shcc-api à l'écoute sur le port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleInstallScript(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("install.sh.tmpl")
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		log.Printf("Erreur template: %v", err)
		return
	}

	data := map[string]string{
		"Host": r.Host,
	}

	w.Header().Set("Content-Type", "text/x-shellscript")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Erreur execution template: %v", err)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "shcc-api",
	})
}

func handleShare(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case http.MethodPost:
		var s api.Share
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "JSON invalide", http.StatusBadRequest)
			return
		}

		s.ID = uuid.New()
		s.CreatedAt = time.Now()

		_, err := database.Pool.Exec(ctx, 
			"INSERT INTO shares (id, owner, \"to\", credentials, created_at, expired_at) VALUES ($1, $2, $3, $4, $5, $6)",
			s.ID, s.Owner, s.To, s.Credentials, s.CreatedAt, s.ExpiredAt)
		
		if err != nil {
			log.Printf("Erreur insertion share: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(s)

	case http.MethodGet:
		to := r.URL.Query().Get("to")
		if to == "" {
			http.Error(w, "Paramètre 'to' manquant", http.StatusBadRequest)
			return
		}

		rows, err := database.Pool.Query(ctx, 
			"SELECT id, owner, \"to\", credentials, created_at, expired_at FROM shares WHERE \"to\" = $1 AND expired_at > NOW()", 
			to)
		if err != nil {
			log.Printf("Erreur query shares: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var shares []api.Share
		for rows.Next() {
			var s api.Share
			if err := rows.Scan(&s.ID, &s.Owner, &s.To, &s.Credentials, &s.CreatedAt, &s.ExpiredAt); err != nil {
				continue
			}
			shares = append(shares, s)
		}

		json.NewEncoder(w).Encode(shares)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	switch r.Method {
	case http.MethodPost:
		var u api.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "JSON invalide", http.StatusBadRequest)
			return
		}
// Upsert sur email
_, err := database.Pool.Exec(ctx, 
	`INSERT INTO users (id, name, email, public_key, created_at, updated_at) 
	 VALUES ($1, $2, $3, $4, NOW(), NOW())
	 ON CONFLICT (email) DO UPDATE 
	 SET name = EXCLUDED.name, 
	     public_key = EXCLUDED.public_key,
	     updated_at = NOW()`,
	uuid.New(), u.Name, u.Email, u.PublicKey)

		if err != nil {
			log.Printf("Erreur upsert user: %v", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	case http.MethodGet:
		email := r.URL.Query().Get("email")
		name := r.URL.Query().Get("name")
		
		var u api.User
		var err error
		if email != "" {
			err = database.Pool.QueryRow(ctx, "SELECT id, name, email, public_key FROM users WHERE email = $1", email).
				Scan(&u.ID, &u.Name, &u.Email, &u.PublicKey)
		} else if name != "" {
			err = database.Pool.QueryRow(ctx, "SELECT id, name, email, public_key FROM users WHERE name = $1", name).
				Scan(&u.ID, &u.Name, &u.Email, &u.PublicKey)
		} else {
			http.Error(w, "Email ou nom manquant", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(u)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
