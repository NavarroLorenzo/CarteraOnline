package transfers

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	ExecuteTransfer(userID int64, input CreateTransferOperation) error
}

type CreateTransferOperation struct {
	FromAccountID      int64
	ToAccountID        int64
	Amount             float64
	TransferID         string
	ExpenseTitle       string
	ExpenseDescription string
	IncomeTitle        string
	IncomeDescription  string
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ExecuteTransfer(userID int64, input CreateTransferOperation) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	lockAccountsQuery := `
		SELECT id, is_active
		FROM accounts
		WHERE user_id = $1
			AND (id = $2 OR id = $3)
		FOR UPDATE
	`

	rows, err := tx.Query(context.Background(), lockAccountsQuery, userID, input.FromAccountID, input.ToAccountID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var sourceFound bool
	var destinationFound bool
	var sourceActive bool
	var destinationActive bool

	for rows.Next() {
		var id int64
		var isActive bool
		if err := rows.Scan(&id, &isActive); err != nil {
			return err
		}

		switch id {
		case input.FromAccountID:
			sourceFound = true
			sourceActive = isActive
		case input.ToAccountID:
			destinationFound = true
			destinationActive = isActive
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if !sourceFound {
		return ErrSourceAccountNotFound
	}
	if !destinationFound {
		return ErrDestinationAccountNotFound
	}
	if !sourceActive {
		return ErrSourceAccountInactive
	}
	if !destinationActive {
		return ErrDestinationAccountInactive
	}

	balanceQuery := `
		SELECT COALESCE(SUM(
			CASE
				WHEN type = 'income' THEN amount
				WHEN type = 'expense' THEN -amount
				ELSE 0
			END
		), 0)
		FROM transactions
		WHERE user_id = $1 AND account_id = $2
	`

	var currentBalance float64
	if err := tx.QueryRow(context.Background(), balanceQuery, userID, input.FromAccountID).Scan(&currentBalance); err != nil {
		return err
	}

	if currentBalance < input.Amount {
		return ErrInsufficientFunds
	}

	insertQuery := `
		INSERT INTO transactions (user_id, title, amount, type, account_id, category, description, transfer_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	if _, err := tx.Exec(
		context.Background(),
		insertQuery,
		userID,
		input.ExpenseTitle,
		input.Amount,
		"expense",
		input.FromAccountID,
		"transfer",
		input.ExpenseDescription,
		input.TransferID,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		context.Background(),
		insertQuery,
		userID,
		input.IncomeTitle,
		input.Amount,
		"income",
		input.ToAccountID,
		"transfer",
		input.IncomeDescription,
		input.TransferID,
	); err != nil {
		return err
	}

	return tx.Commit(context.Background())
}
