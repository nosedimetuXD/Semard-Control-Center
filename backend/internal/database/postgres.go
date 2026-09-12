package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// RunMigrations ejecuta todos los scripts *.up.sql en migrationsDir en orden alfanumérico
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	cleanDir := filepath.Clean(migrationsDir)
	entries, err := os.ReadDir(cleanDir)
	if err != nil {
		return fmt.Errorf("no se pudo leer el directorio de migraciones (%s): %w", cleanDir, err)
	}

	var upFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			upFiles = append(upFiles, entry.Name())
		}
	}
	sort.Strings(upFiles)

	for _, file := range upFiles {
		filePath := filepath.Join(cleanDir, file)
		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error leyendo %s: %w", filePath, err)
		}

		log.Printf("Ejecutando migración SQL: %s...", file)
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("error ejecutando migración %s: %w", file, err)
		}
	}

	log.Printf("Todas las migraciones (%d archivos) fueron aplicadas exitosamente.", len(upFiles))
	return nil
}

// RunInitialMigration mantiene compatibilidad hacia atrás
func RunInitialMigration(ctx context.Context, pool *pgxpool.Pool, migrationFilePath string) error {
	return RunMigrations(ctx, pool, filepath.Dir(migrationFilePath))
}
