package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTSecret   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err == nil {
		log.Println("Variables locales cargadas desde .env")
	}

	cfg := &Config{
		AppPort:     getFirstEnv([]string{"PORT", "APP_PORT"}, "8080"),
		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "postgres"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		JWTSecret:   getEnv("JWT_SECRET", "cambia-este-secreto-en-produccion"),
	}

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
