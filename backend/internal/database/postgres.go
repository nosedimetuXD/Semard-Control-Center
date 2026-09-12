package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error al parsear DATABASE_URL: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 15 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error al crear pool de conexiones: %w", err)
	}

	// Verificar conectividad con Ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("no se pudo conectar con PostgreSQL: %w", err)
	}

	log.Println("Conexión exitosa a la base de datos PostgreSQL")
	return pool, nil
}

// RunInitialMigration ejecuta el script up.sql si encuentra el archivo migrations/000001_init_schema.up.sql
func RunInitialMigration(ctx context.Context, pool *pgxpool.Pool, migrationFilePath string) error {
	// Verificar si el archivo de migración existe
	cleanPath := filepath.Clean(migrationFilePath)
	sqlBytes, err := os.ReadFile(cleanPath)
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo de migración (%s): %w", cleanPath, err)
	}

	log.Printf("Ejecutando migración inicial desde: %s...", cleanPath)
	if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
		return fmt.Errorf("error ejecutando migración inicial SQL: %w", err)
	}

	log.Println("Migración inicial aplicada exitosamente en PostgreSQL")
	return nil
}
