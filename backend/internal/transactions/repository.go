package transactions

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(input CreateTransactionInput) (Transaction, error)
	GetAll(filters TransactionFilters) ([]Transaction, error)
	GetBalance() (float64, error)
	GetBalanceByAccount() (map[int64]float64, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(input CreateTransactionInput) (Transaction, error) {
	var t Transaction

	query := `
		INSERT INTO transactions (title, amount, type, account_id, category, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, amount, type, account_id, category, description, created_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		input.Title,
		input.Amount,
		input.Type,
		input.AccountID,
		input.Category,
		input.Description,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Amount,
		&t.Type,
		&t.AccountID,
		&t.Category,
		&t.Description,
		&t.CreatedAt,
	)

	return t, err
}

func (r *PostgresRepository) GetAll(filters TransactionFilters) ([]Transaction, error) {
	query := `
		SELECT id, title, amount, type, account_id, category, description, created_at
		FROM transactions
		WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if filters.AccountID != nil {
		query += ` AND account_id = $` + strconv.Itoa(argPos)
		args = append(args, *filters.AccountID)
		argPos++
	}

	if filters.Type != nil {
		query += ` AND type = $` + strconv.Itoa(argPos)
		args = append(args, *filters.Type)
		argPos++
	}

	if filters.Category != nil {
		query += ` AND category = $` + strconv.Itoa(argPos)
		args = append(args, *filters.Category)
		argPos++
	}

	if filters.DateFrom != nil {
		query += ` AND created_at >= $` + strconv.Itoa(argPos)
		args = append(args, *filters.DateFrom)
		argPos++
	}

	if filters.DateTo != nil {
		query += ` AND created_at <= $` + strconv.Itoa(argPos)
		args = append(args, *filters.DateTo)
		argPos++
	}

	query += ` ORDER BY id ASC`

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Amount,
			&t.Type,
			&t.AccountID,
			&t.Category,
			&t.Description,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, rows.Err()
}

func (r *PostgresRepository) GetBalance() (float64, error) {
	query := `
		SELECT COALESCE(SUM(
			CASE
				WHEN type = 'income' THEN amount
				WHEN type = 'expense' THEN -amount
				ELSE 0
			END
		), 0)
		FROM transactions
	`

	var balance float64
	err := r.db.QueryRow(context.Background(), query).Scan(&balance)
	return balance, err
}

func (r *PostgresRepository) GetBalanceByAccount() (map[int64]float64, error) {
	query := `
		SELECT
			account_id,
			COALESCE(SUM(
				CASE
					WHEN type = 'income' THEN amount
					WHEN type = 'expense' THEN -amount
					ELSE 0
				END
			), 0) AS balance
		FROM transactions
		GROUP BY account_id
		ORDER BY account_id
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]float64)

	for rows.Next() {
		var accountID int64
		var balance float64

		if err := rows.Scan(&accountID, &balance); err != nil {
			return nil, err
		}

		result[accountID] = balance
	}

	return result, rows.Err()
}
