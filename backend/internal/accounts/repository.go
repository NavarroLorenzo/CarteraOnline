package accounts

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(input CreateAccountInput) (Account, error)
	GetAll(userID int64) ([]Account, error)
	GetByID(userID, id int64) (Account, bool, error)
	ExistsByID(userID, id int64) (bool, error)
	Update(userID, id int64, input UpdateAccountInput) (Account, bool, error)
	Delete(userID, id int64) (bool, error)
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
		INSERT INTO accounts (user_id, name, type)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, name, type, created_at
	`

	err := r.db.QueryRow(context.Background(), query, input.UserID, input.Name, input.Type).
		Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.CreatedAt)

	return account, err
}

func (r *PostgresRepository) GetAll(userID int64) ([]Account, error) {
	query := `
		SELECT id, user_id, name, type, created_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]Account, 0)

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

func (r *PostgresRepository) GetByID(userID, id int64) (Account, bool, error) {
	query := `
		SELECT id, user_id, name, type, created_at
		FROM accounts
		WHERE id = $1 AND user_id = $2
	`

	var account Account
	err := r.db.QueryRow(context.Background(), query, id, userID).
		Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Account{}, false, nil
		}
		return Account{}, false, err
	}

	return account, true, nil
}

func (r *PostgresRepository) ExistsByID(userID, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1 AND user_id = $2)`

	var exists bool
	err := r.db.QueryRow(context.Background(), query, id, userID).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) Update(userID, id int64, input UpdateAccountInput) (Account, bool, error) {
	query := `
		UPDATE accounts
		SET name = $1,
			type = $2
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, name, type, created_at
	`

	var account Account
	err := r.db.QueryRow(context.Background(), query, input.Name, input.Type, id, userID).
		Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Account{}, false, nil
		}

		return Account{}, false, err
	}

	return account, true, nil
}

func (r *PostgresRepository) Delete(userID, id int64) (bool, error) {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return false, err
	}
	defer tx.Rollback(context.Background())

	deleteTransactionsQuery := `DELETE FROM transactions WHERE account_id = $1 AND user_id = $2`
	if _, err := tx.Exec(context.Background(), deleteTransactionsQuery, id, userID); err != nil {
		return false, err
	}

	deleteAccountQuery := `DELETE FROM accounts WHERE id = $1 AND user_id = $2`
	result, err := tx.Exec(context.Background(), deleteAccountQuery, id, userID)
	if err != nil {
		return false, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
