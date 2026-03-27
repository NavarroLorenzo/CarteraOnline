package transactions

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Amount      float64         `json:"amount"`
	Type        TransactionType `json:"type"`
	AccountID   int64           `json:"account_id"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
}

type CreateTransactionInput struct {
	Title       string          `json:"title" binding:"required"`
	Amount      float64         `json:"amount" binding:"required,gt=0"`
	Type        TransactionType `json:"type" binding:"required,oneof=income expense"`
	AccountID   int64           `json:"account_id" binding:"required"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
}
