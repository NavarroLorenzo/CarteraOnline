package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	DatabaseURL        string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	CORSAllowedOrigins []string
}

func LoadConfig() *Config {
	envBeforeDotenv := captureEnvPresence(
		"DATABASE_URL",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSLMODE",
	)

	if err := godotenv.Load(); err == nil {
		log.Println("Variables locales cargadas desde .env")
	}

	cfg := &Config{
		AppPort:            getFirstEnv([]string{"PORT", "APP_PORT"}, "8080"),
		DatabaseURL:        strings.TrimSpace(os.Getenv("DATABASE_URL")),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "postgres"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "cambia-este-secreto-en-produccion"),
		CORSAllowedOrigins: splitCSVEnv("CORS_ALLOWED_ORIGINS"),
	}

	logDatabaseEnvDiagnostics(envBeforeDotenv)

	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getFirstEnv(keys []string, fallback string) string {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value != "" {
			return value
		}
	}

	return fallback
}

func splitCSVEnv(key string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func captureEnvPresence(keys ...string) map[string]bool {
	result := make(map[string]bool, len(keys))
	for _, key := range keys {
		_, exists := os.LookupEnv(key)
		result[key] = exists
	}

	return result
}

func logDatabaseEnvDiagnostics(envBeforeDotenv map[string]bool) {
	log.Printf(
		"DB env diagnostics: DATABASE_URL=%s DB_HOST=%s DB_NAME=%s DB_USER=%s DB_SSLMODE=%s",
		envDiagnostic("DATABASE_URL", envBeforeDotenv),
		envDiagnostic("DB_HOST", envBeforeDotenv),
		envDiagnostic("DB_NAME", envBeforeDotenv),
		envDiagnostic("DB_USER", envBeforeDotenv),
		envDiagnostic("DB_SSLMODE", envBeforeDotenv),
	)
}

func envDiagnostic(key string, envBeforeDotenv map[string]bool) string {
	if envBeforeDotenv[key] {
		return "present(process env)"
	}

	if _, exists := os.LookupEnv(key); exists {
		return "present(.env)"
	}

	return "missing"
}
