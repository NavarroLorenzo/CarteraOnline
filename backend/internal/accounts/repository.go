package accounts

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(input CreateAccountInput) (Account, error)
	GetAll() ([]Account, error)
	GetByID(id int64) (Account, bool, error)
	ExistsByID(id int64) (bool, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(input CreateAccountInput) (Account, error) {
	var account Account

	query := `
		INSERT INTO accounts (name, type)
		VALUES ($1, $2)
		RETURNING id, name, type, created_at
	`

	err := r.db.QueryRow(context.Background(), query, input.Name, input.Type).
		Scan(&account.ID, &account.Name, &account.Type, &account.CreatedAt)

	return account, err
}

func (r *PostgresRepository) GetAll() ([]Account, error) {
	query := `
		SELECT id, name, type, created_at
		FROM accounts
		ORDER BY id ASC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name, &account.Type, &account.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

func (r *PostgresRepository) GetByID(id int64) (Account, bool, error) {
	query := `
		SELECT id, name, type, created_at
		FROM accounts
		WHERE id = $1
	`

	var account Account
	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&account.ID, &account.Name, &account.Type, &account.CreatedAt)

	if err != nil {
		return Account{}, false, err
	}

	return account, true, nil
}

func (r *PostgresRepository) ExistsByID(id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)`

	var exists bool
	err := r.db.QueryRow(context.Background(), query, id).Scan(&exists)
	return exists, err
}
