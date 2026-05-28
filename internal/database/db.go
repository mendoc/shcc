package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// Connect initialise la connexion à PostgreSQL
func Connect() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL non définie")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("erreur parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("erreur création pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("erreur ping db: %w", err)
	}

	Pool = pool
	return nil
}

// Close ferme le pool de connexion
func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
