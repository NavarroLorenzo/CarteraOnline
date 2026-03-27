package transactions

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(input CreateTransactionInput) (Transaction, error)
	GetAll(filters TransactionFilters) ([]Transaction, error)
	GetByID(id int64) (Transaction, bool, error)
	Update(id int64, input UpdateTransactionInput) (Transaction, error)
	Delete(id int64) error
	DeleteByTransferID(transferID string) error
	GetBalance() (float64, error)
	GetBalanceByAccount() (map[int64]float64, error)
	GetSummary(filters TransactionFilters) (TransactionSummary, error)
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
			INSERT INTO transactions (title, amount, type, account_id, category, description, transfer_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, title, amount, type, account_id, category, description, transfer_id, created_at
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
		input.TransferID,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Amount,
		&t.Type,
		&t.AccountID,
		&t.Category,
		&t.Description,
		&t.TransferID,
		&t.CreatedAt,
	)

	return t, err
}

func (r *PostgresRepository) GetAll(filters TransactionFilters) ([]Transaction, error) {
	query := `
		SELECT id, title, amount, type, account_id, category, description, transfer_id, created_at
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
			&t.TransferID,
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

func (r *PostgresRepository) GetSummary(filters TransactionFilters) (TransactionSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS income_total,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS expense_total,
			COALESCE(SUM(CASE
				WHEN type = 'income' THEN amount
				WHEN type = 'expense' THEN -amount
				ELSE 0
			END), 0) AS net_balance,
			COUNT(*) AS transactions_count
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

	var summary TransactionSummary

	err := r.db.QueryRow(context.Background(), query, args...).Scan(
		&summary.IncomeTotal,
		&summary.ExpenseTotal,
		&summary.NetBalance,
		&summary.TransactionsCount,
	)

	return summary, err
}

func (r *PostgresRepository) GetByID(id int64) (Transaction, bool, error) {
	query := `
		SELECT id, title, amount, type, account_id, category, description, transfer_id, created_at
		FROM transactions
		WHERE id = $1
	`

	var t Transaction

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&t.ID,
		&t.Title,
		&t.Amount,
		&t.Type,
		&t.AccountID,
		&t.Category,
		&t.Description,
		&t.TransferID,
		&t.CreatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return Transaction{}, false, nil
		}
		return Transaction{}, false, err
	}

	return t, true, nil
}

func (r *PostgresRepository) Update(id int64, input UpdateTransactionInput) (Transaction, error) {
	query := `
		UPDATE transactions
		SET title = $1,
			amount = $2,
			type = $3,
			account_id = $4,
			category = $5,
			description = $6
		WHERE id = $7
		RETURNING id, title, amount, type, account_id, category, description, transfer_id, created_at
	`

	var t Transaction

	err := r.db.QueryRow(
		context.Background(),
		query,
		input.Title,
		input.Amount,
		input.Type,
		input.AccountID,
		input.Category,
		input.Description,
		id,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Amount,
		&t.Type,
		&t.AccountID,
		&t.Category,
		&t.Description,
		&t.TransferID,
		&t.CreatedAt,
	)

	return t, err
}

func (r *PostgresRepository) Delete(id int64) error {
	query := `DELETE FROM transactions WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *PostgresRepository) DeleteByTransferID(transferID string) error {
	query := `DELETE FROM transactions WHERE transfer_id = $1`
	_, err := r.db.Exec(context.Background(), query, transferID)
	return err
}
