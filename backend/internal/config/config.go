package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	AllowedOrigins     []string
	StoragePath        string
	AllowedEmailDomain string
}

func Load() *Config {
	// Intenta cargar .env si existe (modo local)
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró archivo .env local, leyendo de variables de entorno del sistema")
	}

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/semard_db?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "semard-secret-key-development-change-in-production")
	googleClientID := getEnv("GOOGLE_CLIENT_ID", "")
	googleClientSecret := getEnv("GOOGLE_CLIENT_SECRET", "")
	googleRedirectURL := getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback")
	originsStr := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	storagePath := getEnv("STORAGE_PATH", "./storage")
	emailDomain := getEnv("ALLOWED_EMAIL_DOMAIN", "@unicartagena.edu.co")

	allowedOrigins := strings.Split(originsStr, ",")
	for i, o := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(o)
	}

	return &Config{
		Port:               port,
		DatabaseURL:        dbURL,
		JWTSecret:          jwtSecret,
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GoogleRedirectURL:  googleRedirectURL,
		AllowedOrigins:     allowedOrigins,
		StoragePath:        storagePath,
		AllowedEmailDomain: emailDomain,
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
