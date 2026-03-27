package transactions

import "time"

type Repository interface {
	Create(input CreateTransactionInput) Transaction
	GetAll() []Transaction
	GetBalance() float64
	GetBalanceByAccount() map[int64]float64
}

type InMemoryRepository struct {
	data   []Transaction
	nextID int64
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data:   []Transaction{},
		nextID: 1,
	}
}

func (r *InMemoryRepository) Create(input CreateTransactionInput) Transaction {
	transaction := Transaction{
		ID:          r.nextID,
		Title:       input.Title,
		Amount:      input.Amount,
		Type:        input.Type,
		AccountID:   input.AccountID,
		Category:    input.Category,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}

	r.data = append(r.data, transaction)
	r.nextID++

	return transaction
}

func (r *InMemoryRepository) GetAll() []Transaction {
	return r.data
}

func (r *InMemoryRepository) GetBalance() float64 {
	var balance float64

	for _, t := range r.data {
		if t.Type == Income {
			balance += t.Amount
		} else if t.Type == Expense {
			balance -= t.Amount
		}
	}

	return balance
}

func (r *InMemoryRepository) GetBalanceByAccount() map[int64]float64 {
	balances := make(map[int64]float64)

	for _, t := range r.data {
		if _, exists := balances[t.AccountID]; !exists {
			balances[t.AccountID] = 0
		}

		if t.Type == Income {
			balances[t.AccountID] += t.Amount
		} else if t.Type == Expense {
			balances[t.AccountID] -= t.Amount
		}
	}

	return balances
}
