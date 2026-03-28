package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cartera-app/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(cfg *config.Config) *pgxpool.Pool {
	dsn, source := buildDSN(cfg)
	log.Printf("Inicializando PostgreSQL usando %s", source)

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

func buildDSN(cfg *config.Config) (string, string) {
	if databaseURL := strings.TrimSpace(cfg.DatabaseURL); databaseURL != "" {
		return databaseURL, "DATABASE_URL"
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	return dsn, "variables DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME/DB_SSLMODE"
}

func ensureSchema(pool *pgxpool.Pool) error {
	statements := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			email_verified_at TIMESTAMPTZ NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS accounts (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			name_normalized TEXT NOT NULL,
			type TEXT NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			deleted_at TIMESTAMPTZ NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS transactions (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			amount DOUBLE PRECISION NOT NULL,
			type TEXT NOT NULL,
			account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
			category TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			transfer_id TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
		`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMPTZ NULL`,
		`ALTER TABLE accounts ADD COLUMN IF NOT EXISTS user_id BIGINT`,
		`ALTER TABLE accounts ADD COLUMN IF NOT EXISTS name_normalized TEXT`,
		`ALTER TABLE accounts ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE`,
		`ALTER TABLE accounts ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ NULL`,
		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS user_id BIGINT`,
		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS transfer_id TEXT`,
		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT ''`,
		`
		UPDATE accounts
		SET name = CASE
			WHEN BTRIM(REGEXP_REPLACE(COALESCE(name, ''), '\s+', ' ', 'g')) = '' THEN 'Cuenta ' || id
			ELSE BTRIM(REGEXP_REPLACE(name, '\s+', ' ', 'g'))
		END
		`,
		`
		UPDATE accounts
		SET type = CASE
			WHEN LOWER(BTRIM(type)) IN ('cash', 'bank', 'virtual_wallet', 'credit_card', 'savings') THEN LOWER(BTRIM(type))
			ELSE 'cash'
		END
		`,
		`
		UPDATE accounts
		SET name_normalized = LOWER(BTRIM(REGEXP_REPLACE(name, '\s+', ' ', 'g')))
		WHERE name_normalized IS NULL OR name_normalized = ''
		`,
		`UPDATE accounts SET is_active = TRUE WHERE is_active IS NULL`,
		`
		WITH duplicates AS (
			SELECT
				id,
				ROW_NUMBER() OVER (
					PARTITION BY user_id, name_normalized
					ORDER BY created_at ASC, id ASC
				) AS row_num
			FROM accounts
			WHERE user_id IS NOT NULL
				AND is_active = TRUE
		)
		UPDATE accounts
		SET is_active = FALSE
		FROM duplicates
		WHERE accounts.id = duplicates.id
			AND duplicates.row_num > 1
		`,
		`
		UPDATE transactions
		SET title = CASE
			WHEN BTRIM(REGEXP_REPLACE(COALESCE(title, ''), '\s+', ' ', 'g')) = '' THEN 'Transacción'
			ELSE BTRIM(REGEXP_REPLACE(title, '\s+', ' ', 'g'))
		END
		`,
		`
		UPDATE transactions
		SET category = CASE
			WHEN BTRIM(REGEXP_REPLACE(COALESCE(category, ''), '\s+', ' ', 'g')) = '' AND transfer_id IS NOT NULL THEN 'transfer'
			WHEN BTRIM(REGEXP_REPLACE(COALESCE(category, ''), '\s+', ' ', 'g')) = '' THEN 'general'
			ELSE BTRIM(REGEXP_REPLACE(category, '\s+', ' ', 'g'))
		END
		`,
		`
		UPDATE transactions
		SET description = BTRIM(REGEXP_REPLACE(COALESCE(description, ''), '\s+', ' ', 'g'))
		`,
		`
		UPDATE transactions
		SET type = LOWER(BTRIM(type))
		`,
		`
		UPDATE transactions t
		SET user_id = a.user_id
		FROM accounts a
		WHERE t.account_id = a.id
			AND (t.user_id IS NULL OR t.user_id <> a.user_id)
		`,
		`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conname = 'accounts_type_check'
			) THEN
				ALTER TABLE accounts
				ADD CONSTRAINT accounts_type_check
				CHECK (type IN ('cash', 'bank', 'virtual_wallet', 'credit_card', 'savings'));
			END IF;
		END $$;
		`,
		`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conname = 'accounts_name_not_blank_check'
			) THEN
				ALTER TABLE accounts
				ADD CONSTRAINT accounts_name_not_blank_check
				CHECK (BTRIM(name) <> '');
			END IF;
		END $$;
		`,
		`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conname = 'transactions_type_check'
			) THEN
				ALTER TABLE transactions
				ADD CONSTRAINT transactions_type_check
				CHECK (type IN ('income', 'expense'));
			END IF;
		END $$;
		`,
		`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conname = 'transactions_amount_positive_check'
			) THEN
				ALTER TABLE transactions
				ADD CONSTRAINT transactions_amount_positive_check
				CHECK (amount > 0);
			END IF;
		END $$;
		`,
		`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conname = 'transactions_title_not_blank_check'
			) THEN
				ALTER TABLE transactions
				ADD CONSTRAINT transactions_title_not_blank_check
				CHECK (BTRIM(title) <> '');
			END IF;
		END $$;
		`,
		`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conname = 'transactions_category_not_blank_check'
			) THEN
				ALTER TABLE transactions
				ADD CONSTRAINT transactions_category_not_blank_check
				CHECK (BTRIM(category) <> '');
			END IF;
		END $$;
		`,
		`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_unique_active_name
			ON accounts (user_id, name_normalized)
			WHERE is_active = TRUE AND user_id IS NOT NULL
		`,
		`CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_accounts_user_active ON accounts(user_id, is_active, id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_transfer_id ON transactions(transfer_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_user_created_at ON transactions(user_id, created_at DESC, id DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_user_account_created_at ON transactions(user_id, account_id, created_at DESC, id DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_user_type_created_at ON transactions(user_id, type, created_at DESC, id DESC)`,
	}

	for _, statement := range statements {
		if _, err := pool.Exec(context.Background(), statement); err != nil {
			return err
		}
	}

	return nil
}
