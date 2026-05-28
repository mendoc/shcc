package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Println("Erreur: DATABASE_URL est vide")
		os.Exit(1)
	}

	fmt.Printf("Tentative de connexion à: %s\n", dsn)
	
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Printf("❌ Erreur de connexion: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var version string
	err = conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		fmt.Printf("❌ Erreur de requête: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Succès ! Version PostgreSQL: %s\n", version)
}
