package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"cartera-app/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(cfg *config.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Error parseando configuración de DB: ", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Error conectando a PostgreSQL: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("No se pudo hacer ping a PostgreSQL: ", err)
	}

	if err := ensureSchema(pool); err != nil {
		log.Fatal("No se pudo preparar el esquema de la base de datos: ", err)
	}

	log.Println("Conexión a PostgreSQL OK")
	return pool
}

func ensureSchema(pool *pgxpool.Pool) error {
	statements := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS accounts (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS transactions (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			amount DOUBLE PRECISION NOT NULL,
			type TEXT NOT NULL,
			account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
			category TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			transfer_id TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
		`,
		`ALTER TABLE accounts ADD COLUMN IF NOT EXISTS user_id BIGINT`,
		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS user_id BIGINT`,
		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS transfer_id TEXT`,
		`CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_transfer_id ON transactions(transfer_id)`,
	}

	for _, statement := range statements {
		if _, err := pool.Exec(context.Background(), statement); err != nil {
			return err
		}
	}

	return nil
}
