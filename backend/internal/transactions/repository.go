package transactions

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(userID int64, input CreateTransactionInput) (Transaction, error)
	GetAll(userID int64, filters TransactionFilters) ([]Transaction, error)
	GetByID(userID, id int64) (Transaction, bool, error)
	Update(userID, id int64, input UpdateTransactionInput) (Transaction, error)
	Delete(userID, id int64) error
	DeleteByTransferID(userID int64, transferID string) error
	GetBalance(userID int64) (float64, error)
	GetBalanceByAccount(userID int64) (map[int64]float64, error)
	GetSummary(userID int64, filters TransactionFilters) (TransactionSummary, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(userID int64, input CreateTransactionInput) (Transaction, error) {
	var t Transaction

	query := `
			INSERT INTO transactions (user_id, title, amount, type, account_id, category, description, transfer_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, user_id, title, amount, type, account_id, category, description, transfer_id, created_at
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		userID,
		input.Title,
		input.Amount,
		input.Type,
		input.AccountID,
		input.Category,
		input.Description,
		input.TransferID,
	).Scan(
		&t.ID,
		&t.UserID,
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

func (r *PostgresRepository) GetAll(userID int64, filters TransactionFilters) ([]Transaction, error) {
	query := `
		SELECT id, user_id, title, amount, type, account_id, category, description, transfer_id, created_at
		FROM transactions
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argPos := 2

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

	query += ` ORDER BY created_at DESC, id DESC`

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)

	for rows.Next() {
		var t Transaction
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
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

func (r *PostgresRepository) GetBalance(userID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(
			CASE
				WHEN type = 'income' THEN amount
				WHEN type = 'expense' THEN -amount
				ELSE 0
			END
		), 0)
		FROM transactions
		WHERE user_id = $1
	`

	var balance float64
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&balance)
	return balance, err
}

func (r *PostgresRepository) GetBalanceByAccount(userID int64) (map[int64]float64, error) {
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
		WHERE user_id = $1
		GROUP BY account_id
		ORDER BY account_id
	`

	rows, err := r.db.Query(context.Background(), query, userID)
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

func (r *PostgresRepository) GetSummary(userID int64, filters TransactionFilters) (TransactionSummary, error) {
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
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argPos := 2

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

func (r *PostgresRepository) GetByID(userID, id int64) (Transaction, bool, error) {
	query := `
		SELECT id, user_id, title, amount, type, account_id, category, description, transfer_id, created_at
		FROM transactions
		WHERE id = $1 AND user_id = $2
	`

	var t Transaction

	err := r.db.QueryRow(context.Background(), query, id, userID).Scan(
		&t.ID,
		&t.UserID,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return Transaction{}, false, nil
		}
		return Transaction{}, false, err
	}

	return t, true, nil
}

func (r *PostgresRepository) Update(userID, id int64, input UpdateTransactionInput) (Transaction, error) {
	query := `
		UPDATE transactions
		SET title = $1,
			amount = $2,
			type = $3,
			account_id = $4,
			category = $5,
			description = $6
		WHERE id = $7 AND user_id = $8
		RETURNING id, user_id, title, amount, type, account_id, category, description, transfer_id, created_at
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
		userID,
	).Scan(
		&t.ID,
		&t.UserID,
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

func (r *PostgresRepository) Delete(userID, id int64) error {
	query := `DELETE FROM transactions WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(context.Background(), query, id, userID)
	return err
}

func (r *PostgresRepository) DeleteByTransferID(userID int64, transferID string) error {
	query := `DELETE FROM transactions WHERE transfer_id = $1 AND user_id = $2`
	_, err := r.db.Exec(context.Background(), query, transferID, userID)
	return err
}
