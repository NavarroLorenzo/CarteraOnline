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
	ExistsActiveByNormalizedName(userID int64, nameNormalized string, excludeID *int64) (bool, error)
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
		INSERT INTO accounts (user_id, name, name_normalized, type, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, name_normalized, type, is_active, created_at, deleted_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		input.UserID,
		input.Name,
		input.NameNormalized,
		input.Type,
		input.IsActive,
	).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.NameNormalized,
		&account.Type,
		&account.IsActive,
		&account.CreatedAt,
		&account.DeletedAt,
	)

	return account, err
}

func (r *PostgresRepository) GetAll(userID int64) ([]Account, error) {
	query := `
		SELECT id, user_id, name, name_normalized, type, is_active, created_at, deleted_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY is_active DESC, LOWER(name) ASC, id ASC
	`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]Account, 0)

	for rows.Next() {
		var account Account
		if err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Name,
			&account.NameNormalized,
			&account.Type,
			&account.IsActive,
			&account.CreatedAt,
			&account.DeletedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

func (r *PostgresRepository) GetByID(userID, id int64) (Account, bool, error) {
	query := `
		SELECT id, user_id, name, name_normalized, type, is_active, created_at, deleted_at
		FROM accounts
		WHERE id = $1 AND user_id = $2
	`

	var account Account
	err := r.db.QueryRow(context.Background(), query, id, userID).
		Scan(
			&account.ID,
			&account.UserID,
			&account.Name,
			&account.NameNormalized,
			&account.Type,
			&account.IsActive,
			&account.CreatedAt,
			&account.DeletedAt,
		)

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

func (r *PostgresRepository) ExistsActiveByNormalizedName(userID int64, nameNormalized string, excludeID *int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM accounts
			WHERE user_id = $1
				AND name_normalized = $2
				AND is_active = TRUE
	`

	args := []interface{}{userID, nameNormalized}
	if excludeID != nil {
		query += ` AND id <> $3`
		args = append(args, *excludeID)
	}

	query += `)`

	var exists bool
	err := r.db.QueryRow(context.Background(), query, args...).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) Update(userID, id int64, input UpdateAccountInput) (Account, bool, error) {
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	query := `
		UPDATE accounts
		SET name = $1,
			name_normalized = $2,
			type = $3,
			is_active = $4
		WHERE id = $5 AND user_id = $6
		RETURNING id, user_id, name, name_normalized, type, is_active, created_at, deleted_at
	`

	var account Account
	err := r.db.QueryRow(
		context.Background(),
		query,
		input.Name,
		input.NameNormalized,
		input.Type,
		isActive,
		id,
		userID,
	).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.NameNormalized,
		&account.Type,
		&account.IsActive,
		&account.CreatedAt,
		&account.DeletedAt,
	)
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

	// Borra las transacciones de la cuenta y también cualquier contraparte
	// de transferencia vinculada por transfer_id para no dejar movimientos
	// huérfanos en otras cuentas.
	deleteTransactionsQuery := `
		DELETE FROM transactions
		WHERE user_id = $1
			AND (
				account_id = $2
				OR transfer_id IN (
					SELECT DISTINCT transfer_id
					FROM transactions
					WHERE user_id = $1
						AND account_id = $2
						AND transfer_id IS NOT NULL
				)
			)
	`
	if _, err := tx.Exec(context.Background(), deleteTransactionsQuery, userID, id); err != nil {
		return false, err
	}

	deleteAccountQuery := `
		DELETE FROM accounts
		WHERE id = $1 AND user_id = $2
	`
	result, err := tx.Exec(context.Background(), deleteAccountQuery, id, userID)
	if err != nil {
		return false, err
	}
	if result.RowsAffected() == 0 {
		return false, nil
	}

	if err := tx.Commit(context.Background()); err != nil {
		return false, err
	}

	return true, nil
}
