package database

import (
	"strings"
	"testing"

	"cenz/backend/internal/config"
)

func TestBuildDSNPrioritizesDatabaseURLAndRequiresSSL(t *testing.T) {
	dsn, source := buildDSN(&config.Config{
		DatabaseURL: "postgres://remote_user:secret@remote-host/renderdb",
		DBHost:      "localhost",
		DBPort:      "5432",
		DBUser:      "postgres",
		DBPassword:  "local",
		DBName:      "localdb",
		DBSSLMode:   "disable",
	})

	if source != "DATABASE_URL" {
		t.Fatalf("expected DATABASE_URL source, got %q", source)
	}

	if !strings.Contains(dsn, "remote-host") {
		t.Fatalf("expected remote DATABASE_URL to be used, got %q", dsn)
	}

	if !strings.Contains(dsn, "sslmode=require") {
		t.Fatalf("expected sslmode=require to be appended, got %q", dsn)
	}

	if strings.Contains(dsn, "localhost") {
		t.Fatalf("expected DB_* fallback to be ignored when DATABASE_URL exists, got %q", dsn)
	}
}

func TestBuildDSNPreservesExistingDatabaseURLSSLMode(t *testing.T) {
	dsn, _ := buildDSN(&config.Config{
		DatabaseURL: "postgres://remote_user:secret@remote-host/renderdb?sslmode=verify-full",
	})

	if strings.Count(dsn, "sslmode=") != 1 {
		t.Fatalf("expected existing sslmode to be preserved without duplication, got %q", dsn)
	}

	if !strings.Contains(dsn, "sslmode=verify-full") {
		t.Fatalf("expected existing sslmode value, got %q", dsn)
	}
}

func TestBuildDSNAddsSSLModeToDatabaseURLWithExistingQueryParams(t *testing.T) {
	dsn, source := buildDSN(&config.Config{
		DatabaseURL: "postgres://remote_user:secret@remote-host/renderdb?connect_timeout=10",
		DBHost:      "localhost",
	})

	if source != "DATABASE_URL" {
		t.Fatalf("expected DATABASE_URL source, got %q", source)
	}

	if !strings.Contains(dsn, "connect_timeout=10") {
		t.Fatalf("expected existing query params to be preserved, got %q", dsn)
	}

	if !strings.Contains(dsn, "sslmode=require") {
		t.Fatalf("expected sslmode=require to be appended, got %q", dsn)
	}

	if strings.Contains(dsn, "localhost") {
		t.Fatalf("expected DB_* fallback to be ignored when DATABASE_URL exists, got %q", dsn)
	}
}

func TestBuildDSNFallsBackToLocalDBVariables(t *testing.T) {
	dsn, source := buildDSN(&config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "local",
		DBName:     "cenz_local",
		DBSSLMode:  "disable",
	})

	if source != "variables DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME/DB_SSLMODE" {
		t.Fatalf("unexpected source: %q", source)
	}

	expected := "postgres://postgres:local@localhost:5432/cenz_local?sslmode=disable"
	if dsn != expected {
		t.Fatalf("expected %q, got %q", expected, dsn)
	}
}

func TestEnsureSSLModeSupportsKeywordConnectionStrings(t *testing.T) {
	dsn := ensureSSLMode("host=remote user=render dbname=app", "require")

	if dsn != "host=remote user=render dbname=app sslmode=require" {
		t.Fatalf("unexpected keyword DSN: %q", dsn)
	}
}

func TestEnsureSSLModePreservesKeywordConnectionStringSSLMode(t *testing.T) {
	dsn := ensureSSLMode("host=remote user=render dbname=app sslmode=verify-full", "require")

	if dsn != "host=remote user=render dbname=app sslmode=verify-full" {
		t.Fatalf("expected existing keyword sslmode to be preserved, got %q", dsn)
	}
}
