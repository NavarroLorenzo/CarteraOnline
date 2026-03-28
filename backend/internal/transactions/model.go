package transactions

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

const InitialBalanceCategory = "initial_balance"
const TransferCategory = "transfer"

type Transaction struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"-"`
	Title       string          `json:"title"`
	Amount      float64         `json:"amount"`
	Type        TransactionType `json:"type"`
	AccountID   int64           `json:"account_id"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	TransferID  *string         `json:"transfer_id,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type CreateTransactionInput struct {
	UserID      int64           `json:"-"`
	Title       string          `json:"title" binding:"required"`
	Amount      float64         `json:"amount" binding:"required,gt=0"`
	Type        TransactionType `json:"type" binding:"required,oneof=income expense"`
	AccountID   int64           `json:"account_id" binding:"required"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	TransferID  *string         `json:"transfer_id,omitempty"`
}

type TransactionFilters struct {
	AccountID *int64
	Type      *TransactionType
	Category  *string
	DateFrom  *time.Time
	DateTo    *time.Time
}

type TransactionSummary struct {
	IncomeTotal       float64 `json:"income_total"`
	ExpenseTotal      float64 `json:"expense_total"`
	NetBalance        float64 `json:"net_balance"`
	TransactionsCount int64   `json:"transactions_count"`
}

type UpdateTransactionInput struct {
	Title       string          `json:"title" binding:"required"`
	Amount      float64         `json:"amount" binding:"required,gt=0"`
	Type        TransactionType `json:"type" binding:"required,oneof=income expense"`
	AccountID   int64           `json:"account_id" binding:"required"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
}
